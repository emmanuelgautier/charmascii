package output_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/emmanuelgautier/charmascii/internal/output"
)

func TestWriteTXT_Basic(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.txt")
	err := output.WriteTXT(path, []string{"hello", "world"}, output.Metadata{})
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "hello\nworld\n", string(data))
}

func TestWriteTXT_StripsANSI(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.txt")
	err := output.WriteTXT(path, []string{"\x1b[31mred\x1b[0m", "plain"}, output.Metadata{})
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "red\nplain\n", string(data))
}

func TestWriteTXT_Empty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.txt")
	err := output.WriteTXT(path, []string{}, output.Metadata{})
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "\n", string(data))
}

func TestWriteTXT_InvalidPath(t *testing.T) {
	err := output.WriteTXT("/nonexistent/dir/out.txt", []string{"x"}, output.Metadata{})
	assert.Error(t, err)
}
