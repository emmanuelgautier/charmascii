package output

import (
	"encoding/json"
	"io"
	"strings"
)

// JSONResult is the structured output for --output json.
type JSONResult struct {
	Success  bool         `json:"success"`
	Plain    string       `json:"plain,omitempty"`
	Styled   string       `json:"styled,omitempty"`
	Metadata JSONMetadata `json:"metadata,omitempty"`
	Error    string       `json:"error,omitempty"`
}

// JSONMetadata carries the render parameters used to produce the output.
type JSONMetadata struct {
	Font   string `json:"font"`
	Border string `json:"border"`
	Width  int    `json:"width"`
}

// WriteJSON writes a successful result as a single JSON line to w.
// Plain is always ANSI-free (from result.Lines); Styled may contain ANSI codes.
func WriteJSON(w io.Writer, lines []string, styled string, meta JSONMetadata) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false) // preserve < > & in ASCII art
	return enc.Encode(JSONResult{
		Success:  true,
		Plain:    strings.Join(lines, "\n"),
		Styled:   styled,
		Metadata: meta,
	})
}

// WriteJSONError writes a failure result as a single JSON line to w.
func WriteJSONError(w io.Writer, msg string) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(JSONResult{Success: false, Error: msg})
}
