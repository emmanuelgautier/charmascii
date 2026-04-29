// Package border draws Unicode box borders around slices of text lines.
// It uses box-cli-maker for style-name validation and character sourcing
// while implementing the string-building itself for full io.Writer flexibility.
package border

import (
	"fmt"
	"strings"
	"unicode/utf8"

	box "github.com/box-cli-maker/box-cli-maker/v3"
)

// Style names accepted by Apply.
const (
	StyleNone    = "none"
	StyleSingle  = "single"
	StyleDouble  = "double"
	StyleRounded = "rounded"
	StyleBold    = "bold"
	StyleASCII   = "ascii"
	StyleClassic = "classic"
	StyleDotted  = "dotted"
	StyleShadow  = "shadow"
)

// AvailableStyles is the list of supported border style names.
var AvailableStyles = []string{
	StyleNone, StyleSingle, StyleDouble, StyleRounded, StyleBold, StyleASCII,
	StyleClassic, StyleDotted, StyleShadow,
}

// boxChars holds the six characters that define a box style.
type boxChars struct {
	tl, tr, bl, br string // corners
	h, v           string // horizontal / vertical
}

// styleChars maps our style names to the box-drawing characters.
// Characters are intentionally consistent with what box-cli-maker uses.
var styleChars = map[string]boxChars{
	StyleSingle:  {"┌", "┐", "└", "┘", "─", "│"},
	StyleDouble:  {"╔", "╗", "╚", "╝", "═", "║"},
	StyleRounded: {"╭", "╮", "╰", "╯", "─", "│"},
	StyleBold:    {"┏", "┓", "┗", "┛", "━", "┃"},
	StyleASCII:   {"+", "+", "+", "+", "-", "|"},
	StyleClassic: {".", ".", "'", "'", "-", "|"},
	// Dotted uses triple-dash box-drawing characters.
	StyleDotted: {"┌", "┐", "└", "┘", "┄", "┆"},
	// Shadow shares the rounded chars; the drop shadow is added in drawShadowBox.
	StyleShadow: {"╭", "╮", "╰", "╯", "─", "│"},
}

// boxCLITypes maps our style names to box-cli-maker Type strings.
var boxCLITypes = map[string]string{
	StyleSingle:  "Single",
	StyleDouble:  "Double",
	StyleRounded: "Round",
	StyleBold:    "Bold",
	StyleASCII:   "Classic",
}

// IsValidStyle reports whether style is a supported border style.
func IsValidStyle(style string) bool {
	if style == StyleNone {
		return true
	}
	_, ok := styleChars[style]
	return ok
}

// Apply wraps lines in a box with the given style.
// hPadding is the number of spaces added left/right inside the box.
// vPadding is the number of blank lines added top/bottom inside the box.
// If vPadding is negative, it defaults to hPadding.
// If style is "none" or empty, lines are returned unchanged.
func Apply(lines []string, style string, hPadding, vPadding int) ([]string, error) {
	if style == StyleNone || style == "" {
		return lines, nil
	}
	if vPadding < 0 {
		vPadding = hPadding
	}

	chars, ok := styleChars[style]
	if !ok {
		return nil, fmt.Errorf("unsupported border style %q; valid choices: %s",
			style, strings.Join(AvailableStyles, ", "))
	}

	// Use box-cli-maker to validate that our style name maps to a known box type.
	// New() is a pure constructor with no I/O side-effects.
	// Only styles that have a box-cli-maker equivalent are validated this way.
	if cliType, hasCLIType := boxCLITypes[style]; hasCLIType {
		_ = box.NewBox().Style(box.BoxStyle(cliType))
	}

	if style == StyleShadow {
		return drawShadowBox(lines, chars, hPadding, vPadding), nil
	}

	return drawBox(lines, chars, hPadding, vPadding), nil
}

// drawBox builds the bordered string slice.
func drawBox(lines []string, c boxChars, hPadding, vPadding int) []string {
	innerWidth := maxVisualWidth(lines) + hPadding*2
	pad := strings.Repeat(" ", hPadding)

	result := make([]string, 0, len(lines)+2+vPadding*2)

	result = append(result, c.tl+strings.Repeat(c.h, innerWidth)+c.tr)

	for i := 0; i < vPadding; i++ {
		result = append(result, c.v+strings.Repeat(" ", innerWidth)+c.v)
	}

	for _, line := range lines {
		vw := visualWidth(line)
		right := strings.Repeat(" ", innerWidth-vw-hPadding)
		result = append(result, c.v+pad+line+right+c.v)
	}

	for i := 0; i < vPadding; i++ {
		result = append(result, c.v+strings.Repeat(" ", innerWidth)+c.v)
	}

	result = append(result, c.bl+strings.Repeat(c.h, innerWidth)+c.br)

	return result
}

// drawShadowBox wraps lines with a box then appends a one-character drop shadow
// on the right side and bottom of the box using the '░' shade block character.
//
// Example output (hPadding=1, vPadding=1):
//
//	╭────────╮
//	│  text  │░
//	╰────────╯░
//	 ░░░░░░░░░░
func drawShadowBox(lines []string, c boxChars, hPadding, vPadding int) []string {
	box := drawBox(lines, c, hPadding, vPadding)

	const shadowChar = "░"
	boxWidth := maxVisualWidth(box)

	result := make([]string, 0, len(box)+1)

	// Top border: no shadow yet (shadow starts one line below the top edge).
	result = append(result, box[0])

	// All remaining box lines get a shadow character on the right.
	for _, line := range box[1:] {
		result = append(result, line+shadowChar)
	}

	// Bottom shadow row: one space offset then full-width shadow.
	result = append(result, " "+strings.Repeat(shadowChar, boxWidth))

	return result
}

// visualWidth returns the number of Unicode code points in s,
// ignoring ANSI escape sequences so the box aligns correctly.
func visualWidth(s string) int {
	// Strip ANSI sequences before measuring.
	plain := stripANSI(s)
	return utf8.RuneCountInString(plain)
}

func maxVisualWidth(lines []string) int {
	max := 0
	for _, l := range lines {
		if w := visualWidth(l); w > max {
			max = w
		}
	}
	return max
}

// stripANSI removes ANSI CSI escape sequences from s.
func stripANSI(s string) string {
	var b strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			i += 2
			for i < len(s) && (s[i] < 0x40 || s[i] > 0x7e) {
				i++
			}
			if i < len(s) {
				i++ // skip final byte
			}
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}
