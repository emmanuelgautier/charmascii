package output

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	pngFontSize = 14.0
	pngDPI      = 72.0
	pngPadding  = 20
)

// WritePNG renders lines as a PNG image and writes it to path.
// bgColorName and fgColorName accept "black", "white", or a "#RRGGBB" hex string.
// If meta has content, tEXt metadata chunks are embedded in the PNG.
func WritePNG(path string, lines []string, bgColorName, fgColorName string, meta Metadata) error {
	face, err := loadGoMonoFace(pngFontSize)
	if err != nil {
		return fmt.Errorf("loading font: %w", err)
	}
	defer face.Close()

	metrics := face.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	advance := measureMaxAdvance(face, lines)

	imgW := advance + pngPadding*2
	imgH := lineHeight*len(lines) + pngPadding*2
	if imgW <= pngPadding*2 {
		imgW = 1
	}
	if imgH <= pngPadding*2 {
		imgH = 1
	}

	bg := parseRGBA(bgColorName, color.RGBA{0, 0, 0, 255})
	fg := parseRGBA(fgColorName, color.RGBA{255, 255, 255, 255})

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(img, img.Bounds(), image.NewUniform(bg), image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(fg),
		Face: face,
	}

	for i, line := range lines {
		plain := stripANSI(line)
		y := pngPadding + i*lineHeight + metrics.Ascent.Ceil()
		d.Dot = fixed.P(pngPadding, y)
		d.DrawString(plain)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}

	pngData := buf.Bytes()
	if meta.hasContent() {
		pngData = injectPNGTextChunks(pngData, meta)
	}

	return os.WriteFile(path, pngData, 0o644)
}

// injectPNGTextChunks inserts tEXt metadata chunks before the IEND chunk.
func injectPNGTextChunks(data []byte, meta Metadata) []byte {
	// The IEND chunk is always the last 12 bytes of a valid PNG.
	iendStart := len(data) - 12

	var out bytes.Buffer
	out.Write(data[:iendStart])

	toolStr := fmt.Sprintf("charmascii %s (%s)", meta.Version, meta.URL)
	writePNGtEXtChunk(&out, "Software", toolStr)
	if meta.Command != "" {
		writePNGtEXtChunk(&out, "Comment", meta.Command)
	}

	out.Write(data[iendStart:])
	return out.Bytes()
}

// writePNGtEXtChunk writes a single PNG tEXt chunk to w.
func writePNGtEXtChunk(w *bytes.Buffer, keyword, text string) {
	chunkType := []byte("tEXt")
	data := append([]byte(keyword), 0) // null separator after keyword
	data = append(data, []byte(text)...)

	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(data)))
	w.Write(lenBuf[:])
	w.Write(chunkType)
	w.Write(data)

	crcVal := crc32.NewIEEE()
	crcVal.Write(chunkType)
	crcVal.Write(data)
	var crcBuf [4]byte
	binary.BigEndian.PutUint32(crcBuf[:], crcVal.Sum32())
	w.Write(crcBuf[:])
}

func loadGoMonoFace(size float64) (font.Face, error) {
	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(f, &opentype.FaceOptions{
		Size: size,
		DPI:  pngDPI,
	})
}

func measureMaxAdvance(face font.Face, lines []string) int {
	max := 0
	for _, line := range lines {
		plain := stripANSI(line)
		w := font.MeasureString(face, plain).Ceil()
		if w > max {
			max = w
		}
	}
	return max
}

// parseRGBA parses a color name or "#RRGGBB" hex string.
func parseRGBA(name string, def color.RGBA) color.RGBA {
	namedColors := map[string]color.RGBA{
		"black":   {0, 0, 0, 255},
		"white":   {255, 255, 255, 255},
		"red":     {255, 0, 0, 255},
		"green":   {0, 204, 0, 255},
		"blue":    {0, 102, 255, 255},
		"cyan":    {0, 204, 255, 255},
		"magenta": {255, 0, 255, 255},
		"yellow":  {255, 255, 0, 255},
	}
	name = strings.ToLower(strings.TrimSpace(name))
	if c, ok := namedColors[name]; ok {
		return c
	}
	if strings.HasPrefix(name, "#") && len(name) == 7 {
		var r, g, b uint8
		if _, err := fmt.Sscanf(name, "#%02x%02x%02x", &r, &g, &b); err == nil {
			return color.RGBA{r, g, b, 255}
		}
	}
	return def
}
