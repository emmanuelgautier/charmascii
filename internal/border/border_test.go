package border_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/emmanuelgautier/charmascii/internal/border"
)

func TestApply_None(t *testing.T) {
	input := []string{"hello", "world"}
	out, err := border.Apply(input, border.StyleNone, 1, -1)
	require.NoError(t, err)
	assert.Equal(t, input, out)
}

func TestApply_Single(t *testing.T) {
	out, err := border.Apply([]string{"hi"}, border.StyleSingle, 1, -1)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	// First and last rows must contain corner characters.
	assert.Contains(t, out[0], "┌")
	assert.Contains(t, out[0], "┐")
	assert.Contains(t, out[len(out)-1], "└")
	assert.Contains(t, out[len(out)-1], "┘")
}

func TestApply_Double(t *testing.T) {
	out, err := border.Apply([]string{"hi"}, border.StyleDouble, 1, -1)
	require.NoError(t, err)
	assert.Contains(t, out[0], "╔")
	assert.Contains(t, out[len(out)-1], "╚")
}

func TestApply_Rounded(t *testing.T) {
	out, err := border.Apply([]string{"hi"}, border.StyleRounded, 1, -1)
	require.NoError(t, err)
	assert.Contains(t, out[0], "╭")
}

func TestApply_Bold(t *testing.T) {
	out, err := border.Apply([]string{"hi"}, border.StyleBold, 1, -1)
	require.NoError(t, err)
	assert.Contains(t, out[0], "┏")
}

func TestApply_ASCII(t *testing.T) {
	out, err := border.Apply([]string{"hi"}, border.StyleASCII, 0, -1)
	require.NoError(t, err)
	assert.Contains(t, out[0], "+")
	assert.Contains(t, out[0], "-")
}

func TestApply_InvalidStyle(t *testing.T) {
	_, err := border.Apply([]string{"hi"}, "bogus", 1, -1)
	assert.Error(t, err)
}

func TestApply_MultiLine(t *testing.T) {
	lines := []string{"short", "a much longer line", "mid"}
	out, err := border.Apply(lines, border.StyleSingle, 1, -1)
	require.NoError(t, err)
	// All rows must have the same width.
	width := len([]rune(out[0]))
	for _, row := range out {
		assert.Equal(t, width, len([]rune(row)), "row widths must be equal: %q", row)
	}
}

func TestApply_EmptyString(t *testing.T) {
	out, err := border.Apply([]string{}, "none", 1, -1)
	require.NoError(t, err)
	assert.Empty(t, out)
}

func TestApply_VerticalBars(t *testing.T) {
	out, err := border.Apply([]string{"text"}, border.StyleSingle, 1, -1)
	require.NoError(t, err)
	// Middle rows must start and end with vertical bar.
	for _, row := range out[1 : len(out)-1] {
		assert.True(t, strings.HasPrefix(row, "│"), "row should start with │: %q", row)
		assert.True(t, strings.HasSuffix(row, "│"), "row should end with │: %q", row)
	}
}

func TestIsValidStyle(t *testing.T) {
	for _, s := range border.AvailableStyles {
		assert.True(t, border.IsValidStyle(s), "should be valid: %s", s)
	}
	assert.False(t, border.IsValidStyle("bogus"))
}
