package charmascii_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/emmanuelgautier/charmascii"
	"github.com/emmanuelgautier/charmascii/internal/output"
)

var update = flag.Bool("update", false, "regenerate golden test files")

// TestGenerate_Default checks the basic generation pipeline.
func TestGenerate_Default(t *testing.T) {
	result, err := charmascii.Generate("Hi", charmascii.DefaultOptions())
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Lines)
	assert.NotEmpty(t, result.Styled)
}

func TestGenerate_EmptyText(t *testing.T) {
	result, err := charmascii.Generate("", charmascii.DefaultOptions())
	require.NoError(t, err)
	assert.Empty(t, result.Lines)
}

func TestGenerate_WithColor(t *testing.T) {
	opts := charmascii.DefaultOptions()
	opts.Color = "cyan"
	result, err := charmascii.Generate("A", opts)
	require.NoError(t, err)
	// Styled should contain ANSI codes; Lines should be plain.
	for _, l := range result.Lines {
		assert.NotContains(t, l, "\x1b", "Lines should be plain (no ANSI)")
	}
}

func TestGenerate_NoColor(t *testing.T) {
	opts := charmascii.DefaultOptions()
	opts.Color = "red"
	opts.NoColor = true
	result, err := charmascii.Generate("B", opts)
	require.NoError(t, err)
	assert.NotContains(t, result.Styled, "\x1b", "Styled should be plain when NoColor=true")
}

func TestGenerate_WithBorder(t *testing.T) {
	opts := charmascii.DefaultOptions()
	opts.Border = "single"
	result, err := charmascii.Generate("X", opts)
	require.NoError(t, err)
	assert.NotEmpty(t, result.Lines)
	// Top border line must contain corner char.
	found := false
	for _, l := range result.Lines {
		if strings.Contains(l, "┌") {
			found = true
			break
		}
	}
	assert.True(t, found, "border should contain ┌")
}

func TestGenerate_WithGradient(t *testing.T) {
	opts := charmascii.DefaultOptions()
	opts.Gradient = "red:blue"
	result, err := charmascii.Generate("Go", opts)
	require.NoError(t, err)
	assert.NotEmpty(t, result.Lines)
}

func TestGenerate_InvalidGradient(t *testing.T) {
	opts := charmascii.DefaultOptions()
	opts.Gradient = "nocolon"
	_, err := charmascii.Generate("X", opts)
	assert.Error(t, err)
}

func TestGenerate_InvalidFont(t *testing.T) {
	opts := charmascii.DefaultOptions()
	opts.Font = "nosuchfont"
	_, err := charmascii.Generate("X", opts)
	assert.Error(t, err)
}

func TestGenerate_InvalidBorder(t *testing.T) {
	opts := charmascii.DefaultOptions()
	opts.Border = "nosuchborder"
	_, err := charmascii.Generate("X", opts)
	assert.Error(t, err)
}

func TestGenerate_CenterAlign(t *testing.T) {
	opts := charmascii.DefaultOptions()
	opts.Align = "center"
	opts.Width = 80
	result, err := charmascii.Generate("Hi", opts)
	require.NoError(t, err)
	// At least one line should have leading spaces due to centering.
	hasLeadingSpace := false
	for _, l := range result.Lines {
		if strings.HasPrefix(l, " ") {
			hasLeadingSpace = true
			break
		}
	}
	assert.True(t, hasLeadingSpace, "center-aligned lines should have leading spaces")
}

func TestGenerate_WidthTruncation(t *testing.T) {
	opts := charmascii.DefaultOptions()
	opts.Width = 5
	result, err := charmascii.Generate("Hello World", opts)
	require.NoError(t, err)
	for _, l := range result.Lines {
		assert.LessOrEqual(t, len([]rune(l)), 5, "line should be truncated to width 5")
	}
}

func TestGenerate_AllFonts(t *testing.T) {
	for _, font := range charmascii.ListFonts() {
		t.Run(font, func(t *testing.T) {
			opts := charmascii.DefaultOptions()
			opts.Font = font
			result, err := charmascii.Generate("A", opts)
			require.NoError(t, err, "font: %s", font)
			assert.NotEmpty(t, result.Lines, "font: %s", font)
		})
	}
}

// Golden file tests ─────────────────────────────────────────────────────────

func TestGolden_Standard(t *testing.T) {
	opts := charmascii.DefaultOptions()
	opts.Font = "standard"
	result, err := charmascii.Generate("Hello", opts)
	require.NoError(t, err)

	got := strings.Join(result.Lines, "\n") + "\n"
	checkGolden(t, "testdata/hello_standard.txt", got)
}

func TestGolden_Doom(t *testing.T) {
	opts := charmascii.DefaultOptions()
	opts.Font = "doom"
	result, err := charmascii.Generate("Hello", opts)
	require.NoError(t, err)

	got := strings.Join(result.Lines, "\n") + "\n"
	checkGolden(t, "testdata/hello_doom.txt", got)
}

func TestGolden_SVG(t *testing.T) {
	opts := charmascii.DefaultOptions()
	opts.Font = "banner"
	result, err := charmascii.Generate("Hello", opts)
	require.NoError(t, err)

	// Build the SVG content so we can golden-test it.
	tmpSVG := filepath.Join(t.TempDir(), "banner.svg")
	require.NoError(t, writeSVGGolden(tmpSVG, result.Lines))
	data, err := os.ReadFile(tmpSVG)
	require.NoError(t, err)
	checkGolden(t, "testdata/hello_banner.svg", string(data))
}

func writeSVGGolden(path string, lines []string) error {
	return output.WriteSVG(path, lines, "black", "white", output.Metadata{})
}

// checkGolden compares got against the golden file at path.
// When -update is passed it writes got to path instead.
func checkGolden(t *testing.T, path, got string) {
	t.Helper()
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		return
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Fatalf("golden file %s missing; run with -update to create it", path)
	}
	require.NoError(t, err)
	assert.Equal(t, string(data), got)
}

func TestListFonts(t *testing.T) {
	fonts := charmascii.ListFonts()
	assert.NotEmpty(t, fonts)
	assert.Contains(t, fonts, "standard")
	assert.Contains(t, fonts, "doom")
}

func TestListBorderStyles(t *testing.T) {
	styles := charmascii.ListBorderStyles()
	assert.Contains(t, styles, "none")
	assert.Contains(t, styles, "double")
}

func TestListColors(t *testing.T) {
	colors := charmascii.ListColors()
	assert.Contains(t, colors, "default")
	assert.Contains(t, colors, "cyan")
}
