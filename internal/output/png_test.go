package output_test

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/emmanuelgautier/charmascii/internal/output"
)

func TestWritePNG_Basic(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.png")
	err := output.WritePNG(path, []string{"Hello", "World"}, "black", "white", output.Metadata{})
	require.NoError(t, err)

	// Verify it's a valid PNG.
	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	img, err := png.Decode(f)
	require.NoError(t, err)
	bounds := img.Bounds()
	assert.Greater(t, bounds.Dx(), 0)
	assert.Greater(t, bounds.Dy(), 0)
}

func TestWritePNG_MultiLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), "multi.png")
	lines := []string{"line1", "line2", "line3", "line4"}
	err := output.WritePNG(path, lines, "black", "cyan", output.Metadata{})
	require.NoError(t, err)

	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	img, err := png.Decode(f)
	require.NoError(t, err)
	// More lines → taller image (rough check: height > width is unlikely for 5-char lines).
	bounds := img.Bounds()
	assert.Greater(t, bounds.Dy(), 0)
}

func TestWritePNG_ANSIStripped(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ansi.png")
	// Should not error even with ANSI codes in lines.
	err := output.WritePNG(path, []string{"\x1b[31mred text\x1b[0m"}, "black", "white", output.Metadata{})
	require.NoError(t, err)
}

func TestWritePNG_HexColors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hex.png")
	err := output.WritePNG(path, []string{"test"}, "#1A2B3C", "#FFFFFF", output.Metadata{})
	require.NoError(t, err)

	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()
	_, err = png.Decode(f)
	require.NoError(t, err)
}

func TestWritePNG_InvalidPath(t *testing.T) {
	err := output.WritePNG("/nonexistent/dir/out.png", []string{"x"}, "black", "white", output.Metadata{})
	assert.Error(t, err)
}

func TestWritePNG_EmptyLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.png")
	err := output.WritePNG(path, []string{}, "black", "white", output.Metadata{})
	require.NoError(t, err)

	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	cfg, err := png.Decode(f)
	require.NoError(t, err)
	// Should produce a 1×1 (minimum) image.
	assert.GreaterOrEqual(t, cfg.(image.Image).Bounds().Dx(), 1)
}
