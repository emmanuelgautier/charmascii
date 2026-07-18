# charmascii — full reference

## All flags

| Flag | Default | Values / notes |
|------|---------|-----------------|
| `--font` | `standard` | `standard, big, doom, isometric1, slant, block, 3-d, shadow, banner, bulbhead` (and sometimes more — run `--list-fonts` to get the live list) |
| `--border` | `none` | `none, single, double, rounded, bold, ascii` |
| `--color` | `default` | `default, red, green, blue, cyan, magenta, yellow, white`, or hex `#RRGGBB` / `#RGB` |
| `--border-color` | `default` | same choices as `--color` |
| `--align` | `left` | `left, center, right` |
| `--padding` | `1` | inner horizontal padding inside the border box |
| `--v-padding` | same as `--padding` | inner vertical padding (blank lines) inside the border box; `-1` means "mirror `--padding`" |
| `--output` | `terminal` | `terminal, txt, png, svg, json` |
| `--out-file` | `./output.<ext>` | output file path; ignored for `terminal`/`json` (those go to stdout) |
| `--width` | terminal width | max width in characters; longer lines are truncated |
| `--gradient` | (none) | two colors separated by `:`, e.g. `"red:blue"`, `"#FF0000:#0000FF"` — overrides `--color` when set |
| `--bg-color` | `black` | background color for `png`/`svg` output only |
| `--no-color` | `false` | strip all ANSI codes; auto-enabled when stdout isn't a TTY |
| `--text-shadow` | `false` | adds a `░` drop shadow behind the letters |
| `--list-fonts` | — | prints available fonts and exits |
| `--version` | — | prints version/commit/build date and exits |

## JSON output shape

```json
{
  "success": true,
  "plain": "...",
  "styled": "...",
  "metadata": { "font": "standard", "border": "none", "width": 0 }
}
```

- `plain` — always ANSI-free, safe for any downstream consumer.
- `styled` — may contain ANSI escape codes.
- On failure: `{"success": false, "error": "message"}`.

## MCP server

`charmascii mcp` starts a stdio MCP server exposing a `generate_ascii` tool
— use this instead of shelling out when working inside an MCP-compatible
agent that can register the server directly.

```json
{
  "mcpServers": {
    "charmascii": {
      "command": "charmascii",
      "args": ["mcp"]
    }
  }
}
```

`generate_ascii` accepts: `text` (required), `font`, `border`, `color`,
`align`, `padding`, `width`, `gradient`, `text_shadow`. Returns plain text
(no ANSI) — there's no `output`/`out-file` equivalent over MCP, so it's
terminal-only rendering, not file generation.

## Behavior notes

- Text is a single positional arg — multi-word text needs quotes.
- `--width` truncates rather than wraps; there's no automatic line-wrapping
  of long input text.
- PNG/SVG contrast depends on `--bg-color` vs `--color` — check they're not
  both dark or both light before generating an image the user can't read.
- Gradients need two colors; a single `--gradient "red"` (no colon) is
  invalid.

## PNG rendering limits (verified by generating and inspecting samples)

`--output png` rasterizes text through a bundled Go Mono font with limited
Unicode glyph coverage. This makes some font/border choices that look great
in a terminal come out broken or illegible as an image — the failure isn't
visible until you actually render one, so don't assume a terminal-tested
combo carries over.

**Fonts — legible in PNG:** `ansi_shadow`, `banner`, `3-d`, `slant`.

**Fonts — broken/illegible in PNG:** `doom`, `block`, `standard`, `big`,
`shadow`, `bulbhead`, `isometric1`. These render fine in the terminal; in
PNG their strokes overlap, thin out, or (for `block` and `isometric1`)
scramble into unreadable fragments.

**Borders — legible in PNG:** `none`, `single`, `double`, `rounded`,
`ascii`, `classic`, `shadow`.

**Borders — broken in PNG:** `bold`, `dotted` — their box-drawing
characters (`┏┓┗┛━┃` and `┄┆`) aren't in the bundled font and render as
empty tofu (☐) boxes instead of lines.

SVG output uses `<text>` elements with a CSS font stack instead of
rasterizing through the bundled font, so it likely doesn't share this exact
limitation — but it wasn't verified here, and rendering depends on fonts
available to whatever views the SVG. When in doubt, use the same safe
font/border list for SVG as for PNG.
