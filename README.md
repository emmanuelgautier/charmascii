# charmascii

[![CI](https://github.com/emmanuelgautier/charmascii/actions/workflows/ci.yml/badge.svg)](https://github.com/emmanuelgautier/charmascii/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/emmanuelgautier/charmascii?style=for-the-badge)](https://goreportcard.com/report/github.com/emmanuelgautier/charmascii)
[![Go Reference](https://pkg.go.dev/badge/github.com/emmanuelgautier/charmascii.svg)](https://pkg.go.dev/github.com/emmanuelgautier/charmascii)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**charmascii** converts text to styled ASCII art. It works as both a standalone CLI tool and an importable Go library.

## Features

- 10 FIGlet fonts (standard, doom, slant, banner, and more)
- Unicode box borders (single, double, rounded, bold, ASCII)
- ANSI colors and two-color gradients
- Four output formats: **terminal**, **txt**, **png**, **svg**
- Auto-detects TTY; strips ANSI codes in pipes/CI

## Installation

### Homebrew (macOS/Linux)

```bash
brew install emmanuelgautier/tap/charmascii
```

### Snap (Linux)

```bash
sudo snap install charmascii
```

### Chocolatey (Windows)

```powershell
choco install charmascii
```

### Docker

```bash
# Run directly
docker run --rm ghcr.io/emmanuelgautier/charmascii "Hello World"

# With flags
docker run --rm ghcr.io/emmanuelgautier/charmascii "Hello World" --font doom --border double --color cyan
```

### Go install

```bash
go install github.com/emmanuelgautier/charmascii/cmd/charmascii@latest
```

### Download a release binary

Pre-built binaries for Linux, macOS, and Windows are available on the [releases page](https://github.com/emmanuelgautier/charmascii/releases). Release artifacts are signed with [cosign](https://github.com/sigstore/cosign) — verify with:

```bash
cosign verify-blob \
  --certificate charmascii_<version>_checksums.txt.pem \
  --signature charmascii_<version>_checksums.txt.sig \
  --certificate-identity-regexp "https://github.com/emmanuelgautier/charmascii" \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  charmascii_<version>_checksums.txt
```

## CLI Usage

```bash
charmascii "Hello World"
charmascii "Hello World" --font doom --border double --color cyan
charmascii "v2.0"       --font slant --border rounded --output png --out-file banner.png
charmascii "API"        --font big   --gradient "blue:cyan" --border bold --align center
charmascii --list-fonts
```

### All flags

| Flag | Default | Description |
|------|---------|-------------|
| `--font` | `standard` | FIGlet font: standard, big, doom, isometric1, slant, block, 3d, shadow, banner, bulbhead |
| `--border` | `none` | Border style: none, single, double, rounded, bold, ascii |
| `--color` | `default` | Text color: red, green, blue, cyan, magenta, yellow, white |
| `--border-color` | `default` | Border color (same choices as `--color`) |
| `--align` | `left` | Text alignment: left, center, right |
| `--padding` | `1` | Inner padding inside border box |
| `--output` | `terminal` | Output format: terminal, txt, png, svg |
| `--out-file` | `./output.<ext>` | Output file path |
| `--width` | terminal width | Max width in characters |
| `--gradient` | | Two-color gradient e.g. `"red:blue"` |
| `--bg-color` | `black` | Background color for PNG/SVG output |
| `--no-color` | `false` | Strip all ANSI codes |
| `--list-fonts` | | Print available fonts and exit |
| `--version` | | Print version, commit, and build date |

## Library Usage

```go
import "github.com/emmanuelgautier/charmascii"

opts := charmascii.DefaultOptions()
opts.Font   = "doom"
opts.Border = "double"
opts.Color  = "cyan"

result, err := charmascii.Generate("Hello", opts)
if err != nil {
    log.Fatal(err)
}

fmt.Println(result.Styled)          // ANSI-styled for terminal
fmt.Println(result.Lines[0])        // plain-text first line
```

## Development

```bash
git clone https://github.com/emmanuelgautier/charmascii
cd charmascii

make build        # compile binary → bin/charmascii
make test         # run all tests
make test-race    # run with race detector
make coverage     # generate coverage.html
make lint         # golangci-lint
make test-update  # regenerate golden test files
make snapshot     # local GoReleaser snapshot
```

## Project structure

```
charmascii/
├── charmascii.go              # Public library API
├── cmd/charmascii/main.go     # CLI entry point (Cobra)
├── internal/
│   ├── renderer/            # FIGlet rendering via go-figure
│   ├── border/              # Box drawing via box-cli-maker
│   ├── color/               # ANSI color + gradient via lipgloss
│   └── output/              # terminal / txt / png / svg writers
└── testdata/                # Golden files
```

## License

MIT © [Emmanuel Gautier](https://www.emmanuelgautier.com/) — see [LICENSE](https://github.com/emmanuelgautier/charmascii/blob/main/LICENSE) for details.
