// Package renderer converts text to ASCII art using FIGlet fonts via an
// embedded FIGlet parser, plus a built-in emoji glyph map.
package renderer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// AvailableFonts lists all supported font names.
var AvailableFonts = []string{
	"standard", "big", "doom", "isometric1", "slant",
	"block", "3-d", "shadow", "banner", "bulbhead", "ansi_shadow",
}

// IsValidFont reports whether name is a supported font.
func IsValidFont(name string) bool {
	for _, f := range AvailableFonts {
		if f == name {
			return true
		}
	}
	return false
}

// Render converts text into ASCII art lines using the named font.
// An empty font name defaults to "standard".
// Non-ASCII emoji characters are rendered using the built-in emoji glyph map.
func Render(text, font string) ([]string, error) {
	if font == "" {
		font = "standard"
	}
	if !IsValidFont(font) {
		return nil, fmt.Errorf("unsupported font %q; run --list-fonts for choices", font)
	}

	var lines []string

	if font == "ansi_shadow" {
		lines = renderAnsiShadow(text)
	} else {
		ff, err := loadFontByName(font)
		if err != nil {
			return nil, err
		}
		lines = renderFIGText(text, ff)
	}

	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	maxLen := 0
	for _, l := range lines {
		if w := utf8.RuneCountInString(l); w > maxLen {
			maxLen = w
		}
	}
	for i, l := range lines {
		if w := utf8.RuneCountInString(l); w < maxLen {
			lines[i] = l + strings.Repeat(" ", maxLen-w)
		}
	}

	return lines, nil
}

// renderFIGText renders text character-by-character using the given FIGlet
// font, interleaving emoji glyphs for non-ASCII codepoints.
func renderFIGText(text string, font *figFont) []string {
	rows := make([]string, font.height)

	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		var glyph []string
		if isEmojiRune(ch) {
			// Consume continuation codepoints (variation selectors, ZWJ, …).
			for i+1 < len(runes) && isEmojiContinuation(runes[i+1]) {
				i++
			}
			glyph = getEmojiGlyph(ch, font.height)
		} else {
			glyph = font.getGlyph(ch)
		}

		for row, line := range glyph {
			if row < font.height {
				rows[row] += line
			}
		}
	}

	for i, row := range rows {
		rows[i] = strings.TrimRight(row, " ")
	}

	return rows
}
