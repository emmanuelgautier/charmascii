package color_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/emmanuelgautier/charmascii/internal/color"
)

func TestApplyColor_Default(t *testing.T) {
	lines := []string{"hello", "world"}
	out := color.ApplyColor(lines, "default")
	assert.Equal(t, lines, out)
}

func TestApplyColor_Empty(t *testing.T) {
	lines := []string{"hello"}
	out := color.ApplyColor(lines, "")
	assert.Equal(t, lines, out)
}

func TestApplyColor_Red(t *testing.T) {
	lines := []string{"hello"}
	out := color.ApplyColor(lines, "red")
	// Output should differ from input (ANSI codes added) or be unchanged in no-color envs.
	assert.Len(t, out, 1)
	// Strip ANSI and verify text is preserved.
	assert.Equal(t, "hello", color.StripANSI(out[0]))
}

func TestApplyColor_AllColors(t *testing.T) {
	colors := []string{"red", "green", "blue", "cyan", "magenta", "yellow", "white"}
	for _, c := range colors {
		out := color.ApplyColor([]string{"test"}, c)
		require.Len(t, out, 1, "color: %s", c)
		assert.Equal(t, "test", color.StripANSI(out[0]), "color: %s", c)
	}
}

func TestApplyGradient_ValidGradient(t *testing.T) {
	lines := []string{"line1", "line2", "line3"}
	out, err := color.ApplyGradient(lines, "red:blue")
	require.NoError(t, err)
	assert.Len(t, out, 3)
	for i, l := range out {
		assert.Equal(t, lines[i], color.StripANSI(l), "gradient line %d text should be preserved", i)
	}
}

func TestApplyGradient_InvalidFormat(t *testing.T) {
	_, err := color.ApplyGradient([]string{"x"}, "nocolon")
	assert.Error(t, err)
}

func TestApplyGradient_InvalidColor(t *testing.T) {
	_, err := color.ApplyGradient([]string{"x"}, "notacolor:blue")
	assert.Error(t, err)
}

func TestApplyGradient_SingleLine(t *testing.T) {
	out, err := color.ApplyGradient([]string{"only"}, "cyan:magenta")
	require.NoError(t, err)
	assert.Len(t, out, 1)
}

func TestStripANSI(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"\x1b[31mhello\x1b[0m", "hello"},
		{"\x1b[1;32mgreen bold\x1b[0m", "green bold"},
		{"no escape codes", "no escape codes"},
		{"", ""},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, color.StripANSI(tc.input))
	}
}

func TestStripANSILines(t *testing.T) {
	in := []string{"\x1b[31mred\x1b[0m", "plain"}
	out := color.StripANSILines(in)
	assert.Equal(t, []string{"red", "plain"}, out)
}

func TestIsValidColor(t *testing.T) {
	for _, c := range color.AvailableColors {
		assert.True(t, color.IsValidColor(c))
	}
	assert.False(t, color.IsValidColor("purple"))
	assert.False(t, color.IsValidColor(""))
	// Hex colors are valid.
	assert.True(t, color.IsValidColor("#FF0000"))
	assert.True(t, color.IsValidColor("#f00"))
	assert.False(t, color.IsValidColor("#ZZZZZZ"))
	assert.False(t, color.IsValidColor("FF0000"))
}

func TestIsHexColor(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"#FF0000", true},
		{"#00cc00", true},
		{"#F00", true},
		{"#abc", true},
		{"#GGGGGG", false},
		{"FF0000", false},
		{"#FF00", false},    // 4 digits — invalid
		{"#FF00000", false}, // 7 digits — invalid
		{"", false},
		{"red", false},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, color.IsHexColor(tc.input), "input: %q", tc.input)
	}
}

func TestApplyColor_Hex(t *testing.T) {
	lines := []string{"hello"}
	out := color.ApplyColor(lines, "#FF0000")
	assert.Len(t, out, 1)
	assert.Equal(t, "hello", color.StripANSI(out[0]))
}

func TestApplyColor_HexShort(t *testing.T) {
	lines := []string{"hello"}
	out := color.ApplyColor(lines, "#F00")
	assert.Len(t, out, 1)
	assert.Equal(t, "hello", color.StripANSI(out[0]))
}

func TestApplyGradient_HexColors(t *testing.T) {
	lines := []string{"line1", "line2", "line3"}
	out, err := color.ApplyGradient(lines, "#FF0000:#0000FF")
	require.NoError(t, err)
	assert.Len(t, out, 3)
	for i, l := range out {
		assert.Equal(t, lines[i], color.StripANSI(l))
	}
}

func TestApplyGradient_MixedHexAndName(t *testing.T) {
	lines := []string{"line1", "line2"}
	out, err := color.ApplyGradient(lines, "#FF0000:blue")
	require.NoError(t, err)
	assert.Len(t, out, 2)
}

func TestApplyGradient_InvalidHex(t *testing.T) {
	_, err := color.ApplyGradient([]string{"x"}, "#ZZZZZZ:blue")
	assert.Error(t, err)
}

func TestAvailableColors(t *testing.T) {
	assert.Contains(t, color.AvailableColors, "default")
	assert.Contains(t, color.AvailableColors, "red")
	// No duplicates.
	seen := map[string]bool{}
	for _, c := range color.AvailableColors {
		assert.False(t, seen[c], "duplicate color: %s", c)
		seen[c] = true
	}
}

func TestStripANSI_PreservesSpaces(t *testing.T) {
	s := "\x1b[36m  hello world  \x1b[0m"
	assert.Equal(t, "  hello world  ", color.StripANSI(s))
}

func TestStripANSI_MultipleSequences(t *testing.T) {
	s := "\x1b[1m\x1b[31mbold red\x1b[0m\x1b[0m"
	result := color.StripANSI(s)
	assert.Equal(t, "bold red", result)
	assert.False(t, strings.Contains(result, "\x1b"))
}
