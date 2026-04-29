package renderer

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"strings"
	"sync"
)

//go:embed fonts/*.flf
var fontFS embed.FS

const figASCIIOffset = 32                             // ASCII code of space (first printable)
const figASCIIMax = 126                               // ASCII code of ~ (last printable)
const figNumGlyphs = figASCIIMax - figASCIIOffset + 1 // 95

type figFont struct {
	height int
	glyphs [figNumGlyphs][]string
}

// getGlyph returns the rows for the given ASCII character.
// Non-printable or out-of-range characters fall back to '?'.
func (f *figFont) getGlyph(ch rune) []string {
	if ch == ' ' {
		// Space: two spaces per row (matches go-figure behaviour).
		rows := make([]string, f.height)
		for i := range rows {
			rows[i] = "  "
		}
		return rows
	}
	if ch >= figASCIIOffset && ch <= figASCIIMax {
		if g := f.glyphs[ch-figASCIIOffset]; len(g) > 0 {
			return g
		}
	}
	// Fallback to '?'.
	if g := f.glyphs['?'-figASCIIOffset]; len(g) > 0 {
		return g
	}
	// Last-resort: one space per row.
	rows := make([]string, f.height)
	for i := range rows {
		rows[i] = " "
	}
	return rows
}

var (
	fontCacheMu sync.Mutex
	fontCacheM  = map[string]*figFont{}
)

// loadFontByName loads a FIGlet font from the embedded filesystem (cached).
func loadFontByName(name string) (*figFont, error) {
	fontCacheMu.Lock()
	defer fontCacheMu.Unlock()
	if f, ok := fontCacheM[name]; ok {
		return f, nil
	}
	data, err := fontFS.ReadFile("fonts/" + name + ".flf")
	if err != nil {
		return nil, fmt.Errorf("font %q not found", name)
	}
	f, err := parseFIGFont(data)
	if err != nil {
		return nil, fmt.Errorf("parse font %q: %w", name, err)
	}
	fontCacheM[name] = f
	return f, nil
}

// parseFIGFont parses a FIGlet .flf font file.
func parseFIGFont(data []byte) (*figFont, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 1<<20), 1<<20)

	var header string
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "flf2") {
			header = scanner.Text()
			break
		}
	}
	if header == "" {
		return nil, fmt.Errorf("not a FIGlet font (missing flf2 header)")
	}

	if len(header) < 7 {
		return nil, fmt.Errorf("header too short: %q", header)
	}
	hardblank := header[5]

	var height, baseline, maxLen, oldLayout, commentLines int
	n, _ := fmt.Sscanf(header[6:], " %d %d %d %d %d",
		&height, &baseline, &maxLen, &oldLayout, &commentLines)
	if n < 2 {
		return nil, fmt.Errorf("cannot parse height from header %q", header)
	}

	for i := 0; i < commentLines; i++ {
		if !scanner.Scan() {
			break
		}
	}

	font := &figFont{height: height}

	for idx := 0; idx < figNumGlyphs; idx++ {
		glyph, err := readFIGGlyph(scanner, height, hardblank)
		if err != nil {
			return nil, fmt.Errorf("read glyph idx %d (%q): %w", idx, rune(idx+figASCIIOffset), err)
		}
		ch := rune(idx + figASCIIOffset)
		if ch == ' ' {
			for r := range glyph {
				glyph[r] = "  "
			}
		}
		font.glyphs[idx] = glyph
	}

	return font, nil
}

// readFIGGlyph reads one character definition (exactly height rows).
// Each row ends with one or more '@' delimiters which are stripped.
// The hardblank byte is replaced with a regular space.
func readFIGGlyph(scanner *bufio.Scanner, height int, hardblank byte) ([]string, error) {
	rows := make([]string, 0, height)
	for i := 0; i < height; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected EOF at row %d/%d", i, height)
		}
		line := scanner.Text()
		trimmed := strings.TrimRight(line, "@")
		if hardblank != ' ' && hardblank != 0 {
			trimmed = strings.ReplaceAll(trimmed, string(hardblank), " ")
		}
		rows = append(rows, trimmed)
	}
	return rows, nil
}
