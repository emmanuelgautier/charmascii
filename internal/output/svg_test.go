package output_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/emmanuelgautier/charmascii/internal/output"
)

func TestWriteSVG_Basic(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.svg")
	err := output.WriteSVG(path, []string{"Hello", "World"}, "black", "white", output.Metadata{})
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	s := string(data)

	assert.True(t, strings.HasPrefix(s, "<?xml"), "should start with XML declaration")
	assert.Contains(t, s, "<svg ")
	assert.Contains(t, s, "</svg>")
	assert.Contains(t, s, "Hello")
	assert.Contains(t, s, "World")
}

func TestWriteSVG_ContainsFillColor(t *testing.T) {
	path := filepath.Join(t.TempDir(), "color.svg")
	err := output.WriteSVG(path, []string{"hi"}, "navy", "cyan", output.Metadata{})
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	s := string(data)

	assert.Contains(t, s, "navy")
	assert.Contains(t, s, "cyan")
}

func TestWriteSVG_XMLEscape(t *testing.T) {
	path := filepath.Join(t.TempDir(), "escape.svg")
	err := output.WriteSVG(path, []string{"a & b < c > d"}, "black", "white", output.Metadata{})
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	s := string(data)

	assert.Contains(t, s, "&amp;")
	assert.Contains(t, s, "&lt;")
	assert.Contains(t, s, "&gt;")
}

func TestWriteSVG_ANSIStripped(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ansi.svg")
	err := output.WriteSVG(path, []string{"\x1b[36mhi\x1b[0m"}, "black", "white", output.Metadata{})
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	s := string(data)

	assert.NotContains(t, s, "\x1b", "ANSI codes should be stripped from SVG")
	assert.Contains(t, s, ">hi<")
}

func TestWriteSVG_InvalidPath(t *testing.T) {
	err := output.WriteSVG("/nonexistent/dir/out.svg", []string{"x"}, "black", "white", output.Metadata{})
	assert.Error(t, err)
}

func TestWriteSVG_DefaultColors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "default.svg")
	// Empty color strings should fall back to defaults.
	err := output.WriteSVG(path, []string{"test"}, "", "", output.Metadata{})
	require.NoError(t, err)
}

func TestWriteSVG_MultiLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), "multi.svg")
	lines := []string{"line one", "line two", "line three"}
	err := output.WriteSVG(path, lines, "black", "green", output.Metadata{})
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	s := string(data)

	assert.Contains(t, s, "line one")
	assert.Contains(t, s, "line two")
	assert.Contains(t, s, "line three")
}

func TestWriteSVG_XMLSpacePreserve(t *testing.T) {
	path := filepath.Join(t.TempDir(), "space.svg")
	err := output.WriteSVG(path, []string{"  spaced  "}, "black", "white", output.Metadata{})
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), `xml:space="preserve"`)
}
