package output

// Metadata holds tool attribution info embedded in generated files.
type Metadata struct {
	Version string // tool version, e.g. "1.2.3"
	Command string // full CLI invocation, e.g. `charmascii "Hello" --font doom`
	URL     string // project URL, e.g. "https://github.com/emmanuelgautier/charmascii"
}

func (m Metadata) hasContent() bool {
	return m.Version != "" || m.Command != "" || m.URL != ""
}
