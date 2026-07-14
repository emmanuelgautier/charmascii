---
name: charmascii-go
description: Write or edit Go code that imports the charmascii package (github.com/emmanuelgautier/charmascii) to generate styled ASCII art — FIGlet-style banners, unicode borders, ANSI/hex colors, gradients, and drop shadows — programmatically via the Options struct and Generate() function, instead of shelling out to the CLI. Use this whenever the task involves Go source that needs to produce ASCII art banners, startup/splash screens for a Go CLI tool or TUI, embedding a logo/title into a Go program's output, or wiring charmascii into another Go application — even if the user just says "add a banner to my Go CLI", "print a startup logo in my Go program", "use charmascii as a library", or pastes Go code and asks for ASCII art output. If the task is instead about running the `charmascii` binary from a shell or README (not Go source), prefer the companion "charmascii" CLI skill.

---

# charmascii (Go library)

Package `github.com/emmanuelgautier/charmascii` converts text to styled ASCII
art. It's the same engine as the `charmascii` CLI, exposed as a small Go API:
an `Options` struct in, a `Result` struct out. The whole public surface is
one file (`charmascii.go`, ~150 lines of actual logic) — read it directly in
the target module's vendor/module cache if something below seems off, rather
than guessing.

## Quick start

```go
import "github.com/emmanuelgautier/charmascii"

opts := charmascii.DefaultOptions()
opts.Font = "doom"
opts.Border = "double"
opts.Color = "cyan"

result, err := charmascii.Generate("Hello", opts)
if err != nil {
    return fmt.Errorf("ascii render: %w", err)
}

fmt.Println(result.Styled)   // ANSI-styled, ready for terminal stdout
fmt.Println(result.Lines[0]) // plain text, first row — no ANSI codes
```

Always start from `DefaultOptions()` and override only the fields the user
asked for. It's a value type (not a pointer), so mutating the returned
`Options` never affects other call sites.

Add the dependency in a module outside this repo with:

```bash
go get github.com/emmanuelgautier/charmascii@latest
```

If you're editing a `.go` file *inside the charmascii repo itself*, the
package is still imported by its full path (`github.com/emmanuelgautier/charmascii`)
even though you're in the same module — there's no internal shortcut import.

## API surface

```go
func DefaultOptions() Options
func Generate(text string, opts Options) (*Result, error)
func ListFonts() []string
func ListBorderStyles() []string
func ListColors() []string
```

`Options` fields (all optional — zero value falls back sensibly except
where noted):

| Field | Type | Meaning |
|---|---|---|
| `Font` | `string` | FIGlet font name. Default `"standard"`. |
| `Border` | `string` | Border style. Default `"none"`. |
| `Color` | `string` | Text foreground color name or hex. Default `"default"`. |
| `BorderColor` | `string` | Border foreground color. Default `"default"`. |
| `Align` | `string` | `"left"`, `"center"`, or `"right"`. Default `"left"`. |
| `Padding` | `int` | Inner horizontal padding inside a border box. Default `1`. |
| `VPadding` | `int` | Inner vertical padding (blank lines). **Negative means "use `Padding`'s value"** — leave it at `DefaultOptions()`'s `-1` unless the user wants vertical padding to differ from horizontal. |
| `Width` | `int` | Max output width in characters. `0` = no limit. |
| `Gradient` | `string` | Two-color gradient, `"color1:color2"`. Wins over `Color` if both are set. |
| `BgColor` | `string` | Background color, only matters for PNG output. Default `"black"`. |
| `NoColor` | `bool` | Strip ANSI from `Result.Styled`. |
| `TextShadow` | `bool` | Adds a one-character `░` drop shadow behind letters. |

`Result`:

| Field | Type | Meaning |
|---|---|---|
| `Lines` | `[]string` | Plain text, no ANSI, one element per row. Use for txt/png/svg-style output or any non-terminal consumer. |
| `Styled` | `string` | Full ANSI-styled block, newline-joined, ready to `fmt.Println` in a terminal. |

## Don't hardcode font/border/color names

Font, border, and color name lists live in the renderer/border/color
packages and can grow between versions. Call `ListFonts()`,
`ListBorderStyles()`, or `ListColors()` at runtime — e.g. to validate user
input or populate a `--font` flag's help text — instead of copying a list
into code you write. The same caution applies here as in the CLI skill's
flag reference: a hand-copied list can silently drift from what the
installed version actually supports.

## Gotchas learned from the source (not just the doc comments)

- **Empty text is not an error.** `Generate("", opts)` returns
  `&Result{}, nil` — `Lines` will be `nil` and `Styled` `""`. Don't add a
  manual empty-string guard before calling it; the function already handles
  it, and adding your own just duplicates the check.
- **`Gradient` and `Color` are mutually exclusive in effect, not validated.**
  Setting both doesn't error — `Gradient` silently wins. If a user says
  "cyan gradient to blue" (a single color, not two), that's `Color`, not
  `Gradient`; don't force it into the two-color gradient format.
- **`Width` truncates before alignment.** If you set both `Width` and
  `Align: "center"`, lines are cut to `Width` runes first, then centered —
  so a long line and a short `Width` can end up looking unbalanced. That's
  expected behavior, not a bug to work around.
- **`result.Lines` is always ANSI-free**, regardless of `NoColor`. Reach for
  `Lines` (not `Styled` + manual stripping) whenever the destination isn't a
  terminal — writing to a file, embedding in an image, returning JSON from
  an HTTP handler, etc. `NoColor` only affects `Styled`.
- **`Generate` returns `(*Result, error)`** — always check the error before
  touching the result. It surfaces invalid font names, malformed gradient
  strings (wrong `"color1:color2"` shape), and unknown border styles;
  letting a bad `Font`/`Border`/`Gradient` value reach the user as a panic
  or garbled output is avoidable by just checking `err`.

## Examples

**Input:** "give me a Go function that returns a doom-font banner with a double border in cyan"
```go
func Banner(text string) (string, error) {
    opts := charmascii.DefaultOptions()
    opts.Font = "doom"
    opts.Border = "double"
    opts.Color = "cyan"

    result, err := charmascii.Generate(text, opts)
    if err != nil {
        return "", fmt.Errorf("ascii render: %w", err)
    }
    return result.Styled, nil
}
```

**Input:** "I'm building a CLI tool in Go and want to print a startup banner using this library — plain text, no color, since it might run in a pipe"
```go
opts := charmascii.DefaultOptions()
opts.Font = "slant"
opts.NoColor = true

result, err := charmascii.Generate("MyTool", opts)
if err != nil {
    log.Fatalf("ascii render: %v", err)
}
fmt.Println(result.Styled)
```
`NoColor: true` here rather than switching to `result.Lines` — `Styled`
with `NoColor` set still applies alignment/border/shadow layout, `Lines` is
the pre-color-stage plain text (both work for a color-free terminal print,
but `Styled` is the one that matches what `NoColor` is actually for).

**Input:** "write a small Go program that prints one banner per available font, so I can eyeball them"
```go
package main

import (
    "fmt"

    "github.com/emmanuelgautier/charmascii"
)

func main() {
    for _, font := range charmascii.ListFonts() {
        opts := charmascii.DefaultOptions()
        opts.Font = font

        result, err := charmascii.Generate(font, opts)
        if err != nil {
            fmt.Printf("%s: error: %v\n", font, err)
            continue
        }
        fmt.Println(result.Styled)
    }
}
```
Iterating `ListFonts()` instead of a hardcoded slice is the point of this
example — it stays correct as fonts are added or renamed.

**Input:** "I need the ASCII art as JSON-serializable data for an HTTP API, not printed to a terminal"
```go
opts := charmascii.DefaultOptions()
opts.Font = "banner"

result, err := charmascii.Generate(text, opts)
if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
}
json.NewEncoder(w).Encode(map[string]any{
    "lines": result.Lines, // plain, no ANSI — safe to serialize directly
})
```
