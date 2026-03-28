package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/emmanuelgautier/charmascii/internal/output"
)

func TestWriteTerminal_PlainText(t *testing.T) {
	var buf bytes.Buffer
	err := output.WriteTerminal(&buf, []string{"hello", "world"}, true)
	require.NoError(t, err)
	assert.Equal(t, "hello\nworld\n", buf.String())
}

func TestWriteTerminal_StripsANSIWhenNoColor(t *testing.T) {
	var buf bytes.Buffer
	lines := []string{"\x1b[31mred text\x1b[0m", "plain"}
	err := output.WriteTerminal(&buf, lines, true)
	require.NoError(t, err)
	result := buf.String()
	assert.False(t, strings.Contains(result, "\x1b"), "ANSI codes should be stripped")
	assert.Contains(t, result, "red text")
}

func TestWriteTerminal_PreservesANSIWhenColor(t *testing.T) {
	// Write to a plain bytes.Buffer (not a terminal).
	// When noColor=false but writer is not *os.File, the auto-strip doesn't kick in.
	var buf bytes.Buffer
	lines := []string{"\x1b[31mred\x1b[0m"}
	// noColor=false but buf is not *os.File, so no auto-strip.
	err := output.WriteTerminal(&buf, lines, false)
	require.NoError(t, err)
	// The ANSI codes should be preserved since buf is not a *os.File non-TTY.
	assert.Contains(t, buf.String(), "\x1b[31m")
}

func TestWriteTerminal_Empty(t *testing.T) {
	var buf bytes.Buffer
	err := output.WriteTerminal(&buf, []string{}, false)
	require.NoError(t, err)
	assert.Equal(t, "", buf.String())
}
