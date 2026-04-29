// Package charmascii converts text to styled ASCII art.
//
// It can be used both as a standalone CLI tool and as an importable Go library.
//
// Basic usage:
//
//	result, err := charmascii.Generate("Hello", charmascii.DefaultOptions())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.Styled)
package charmascii

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/emmanuelgautier/charmascii/internal/border"
	"github.com/emmanuelgautier/charmascii/internal/color"
	"github.com/emmanuelgautier/charmascii/internal/renderer"
)

// Options configures ASCII art generation.
type Options struct {
	// Font is the FIGlet font name (default: "standard").
	Font string
	// Border is the border style (default: "none").
	Border string
	// Color is the text foreground color (default: "default").
	Color string
	// BorderColor is the border foreground color (default: "default").
	BorderColor string
	// Align is the horizontal alignment: "left", "center", or "right" (default: "left").
	Align string
	// Padding is the inner horizontal padding inside the border box (default: 1).
	Padding int
	// VPadding is the inner vertical padding (blank lines) inside the border box.
	// Negative means use Padding value.
	VPadding int
	// Width limits the maximum output width in characters. 0 means no limit.
	Width int
	// Gradient applies a two-color gradient in "color1:color2" format.
	Gradient string
	// BgColor is the background color for PNG output (default: "black").
	BgColor string
	// NoColor strips all ANSI codes from the output.
	NoColor bool
	// TextShadow adds a one-character drop shadow (░) behind the ASCII-art letters.
	TextShadow bool
}

// DefaultOptions returns Options with production-ready defaults.
func DefaultOptions() Options {
	return Options{
		Font:     "standard",
		Border:   "none",
		Color:    "default",
		Align:    "left",
		Padding:  1,
		VPadding: -1,
		BgColor:  "black",
	}
}

// Result holds the generated ASCII art in both plain and styled forms.
type Result struct {
	// Lines is the plain-text output (no ANSI escape codes), one element per row.
	// These are suitable for txt/png/svg output formats.
	Lines []string
	// Styled is the ANSI-styled output joined by newlines, ready for terminal display.
	Styled string
}

// Generate converts text to ASCII art according to opts.
func Generate(text string, opts Options) (*Result, error) {
	if text == "" {
		return &Result{}, nil
	}

	lines, err := renderer.Render(text, opts.Font)
	if err != nil {
		return nil, err
	}

	if opts.Width > 0 {
		lines = truncateLines(lines, opts.Width)
	}

	if opts.Align != "" && opts.Align != "left" {
		w := maxWidth(lines)
		if opts.Width > 0 && opts.Width > w {
			w = opts.Width
		}
		lines = align(lines, opts.Align, w)
	}

	if opts.TextShadow {
		lines = applyTextShadow(lines)
	}

	styledLines := lines
	switch {
	case opts.Gradient != "":
		styledLines, err = color.ApplyGradient(lines, opts.Gradient)
		if err != nil {
			return nil, fmt.Errorf("gradient: %w", err)
		}
	case opts.Color != "" && opts.Color != "default":
		styledLines = color.ApplyColor(lines, opts.Color)
	}

	if opts.Border != "" && opts.Border != "none" {
		styledLines, err = border.Apply(styledLines, opts.Border, opts.Padding, opts.VPadding)
		if err != nil {
			return nil, fmt.Errorf("border: %w", err)
		}
	}

	plainLines := color.StripANSILines(styledLines)

	styled := strings.Join(styledLines, "\n")
	if opts.NoColor {
		styled = strings.Join(plainLines, "\n")
	}

	return &Result{
		Lines:  plainLines,
		Styled: styled,
	}, nil
}

// ListFonts returns all available FIGlet font names.
func ListFonts() []string {
	return renderer.AvailableFonts
}

// ListBorderStyles returns all available border style names.
func ListBorderStyles() []string {
	return border.AvailableStyles
}

// ListColors returns all available color names.
func ListColors() []string {
	return color.AvailableColors
}

// align adjusts each line to left/center/right within width columns.
func align(lines []string, alignment string, width int) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		vw := utf8.RuneCountInString(line)
		switch alignment {
		case "center":
			pad := (width - vw) / 2
			if pad > 0 {
				result[i] = strings.Repeat(" ", pad) + line
			} else {
				result[i] = line
			}
		case "right":
			pad := width - vw
			if pad > 0 {
				result[i] = strings.Repeat(" ", pad) + line
			} else {
				result[i] = line
			}
		default:
			result[i] = line
		}
	}
	return result
}

// truncateLines trims each line to at most width runes.
func truncateLines(lines []string, width int) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		runes := []rune(line)
		if len(runes) > width {
			result[i] = string(runes[:width])
		} else {
			result[i] = line
		}
	}
	return result
}

// applyTextShadow adds a one-character drop shadow (░) behind the ASCII-art
// letters. The shadow is offset one column to the right and one row below each
// non-space character, and is only visible where the original text is absent.
func applyTextShadow(lines []string) []string {
	const shadowChar = '░'
	const space = ' '

	rows := len(lines)
	if rows == 0 {
		return lines
	}

	// Convert lines to rune slices for O(1) column access.
	runes := make([][]rune, rows)
	for i, l := range lines {
		runes[i] = []rune(l)
	}

	origAt := func(r, c int) rune {
		if r < 0 || r >= rows || c < 0 || c >= len(runes[r]) {
			return space
		}
		return runes[r][c]
	}

	cols := maxWidth(lines) + 1 // +1 for rightmost shadow column
	totalRows := rows + 1       // +1 for bottom shadow row

	result := make([]string, totalRows)
	for r := 0; r < totalRows; r++ {
		row := make([]rune, cols)
		for c := 0; c < cols; c++ {
			if ch := origAt(r, c); ch != space {
				row[c] = ch
			} else if origAt(r-1, c-1) != space {
				row[c] = shadowChar
			} else {
				row[c] = space
			}
		}
		result[r] = strings.TrimRight(string(row), " ")
	}
	return result
}

// maxWidth returns the maximum rune count across lines.
func maxWidth(lines []string) int {
	max := 0
	for _, l := range lines {
		if n := utf8.RuneCountInString(l); n > max {
			max = n
		}
	}
	return max
}
