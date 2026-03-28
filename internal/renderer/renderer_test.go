package renderer_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/emmanuelgautier/charmascii/internal/renderer"
)

func TestRender_Standard(t *testing.T) {
	lines, err := renderer.Render("Hi", "standard")
	require.NoError(t, err)
	assert.NotEmpty(t, lines)
	// ASCII art must be non-trivial (more than one line, non-empty content).
	hasContent := false
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			hasContent = true
			break
		}
	}
	assert.True(t, hasContent, "rendered output should contain non-blank lines")
}

func TestRender_DoomFont(t *testing.T) {
	lines, err := renderer.Render("Hi", "doom")
	require.NoError(t, err)
	assert.NotEmpty(t, lines)
}

func TestRender_DefaultFont(t *testing.T) {
	lines, err := renderer.Render("A", "")
	require.NoError(t, err)
	assert.NotEmpty(t, lines)
}

func TestRender_InvalidFont(t *testing.T) {
	_, err := renderer.Render("X", "nonexistentfont123")
	assert.Error(t, err)
}

func TestRender_NoTrailingBlankLines(t *testing.T) {
	lines, err := renderer.Render("Go", "standard")
	require.NoError(t, err)
	if len(lines) > 0 {
		last := lines[len(lines)-1]
		assert.NotEqual(t, "", strings.TrimSpace(last), "last line should not be blank")
	}
}

func TestIsValidFont(t *testing.T) {
	assert.True(t, renderer.IsValidFont("standard"))
	assert.True(t, renderer.IsValidFont("doom"))
	assert.True(t, renderer.IsValidFont("slant"))
	assert.False(t, renderer.IsValidFont(""))
	assert.False(t, renderer.IsValidFont("notafont"))
}

func TestAvailableFonts(t *testing.T) {
	assert.Contains(t, renderer.AvailableFonts, "standard")
	assert.Contains(t, renderer.AvailableFonts, "doom")
	assert.Len(t, renderer.AvailableFonts, 11)
}
