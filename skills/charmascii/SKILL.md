---
name: charmascii
description: Generate styled ASCII art banners and text art with the charmascii CLI — FIGlet fonts, unicode borders, ANSI/hex colors, two-color gradients, drop shadows, and five output formats (terminal, txt, png, svg, json). Use this whenever the user wants ASCII art, a text banner, a figlet-style header, decorative/big terminal text, a logo or title for a README/CLI splash screen, or a styled text image — even if they never say "charmascii" or "ASCII art" explicitly, e.g. "make my project name look cool in the terminal", "big text banner for my README", "cli splash screen for my app", "figlet my company name", "turn this into block letters". Also covers scripted/agent use via `--output json` or the `charmascii mcp` server.
---

# charmascii

CLI that renders text as styled ASCII art: FIGlet fonts, borders, colors,
gradients, and five output formats. Full flag/font/color reference is in
`reference.md` — read it before crafting a non-trivial command (custom
colors, gradients, image output) rather than guessing flag names.

## Before generating

Confirm the binary is reachable:

```bash
which charmascii || echo "not found"
```

If missing, tell the user how to install it (don't try to build it from
source on their behalf — that assumes a repo checkout they may not have):

```bash
brew install emmanuelgautier/tap/charmascii   # macOS/Linux
go install github.com/emmanuelgautier/charmascii/cmd/charmascii@latest
```

## Core command shape

```bash
charmascii "TEXT" --font FONT --border STYLE --color COLOR --align left|center|right --output FORMAT
```

`TEXT` is a single positional argument — quote it. Every other option is a
flag with a sane default, so start minimal (`charmascii "Hi"`) and add flags
only for what the user actually asked for. Run `charmascii --list-fonts` if
unsure which fonts are installed — the flag help text and the actual font
set can drift (e.g. `ansi_shadow` is available but not always documented).

## Picking style from a vague request

Users rarely name flags directly ("make it look bold and cool" — not
"--font doom --border double"). Translate intent:

| User says... | Reach for... |
|---|---|
| "bold", "loud", "impact" | `--font doom` or `--font block`, `--border bold` |
| "clean", "modern", "minimal" | `--font standard` or `--font slant`, `--border single` or no border |
| "retro", "3D", "chunky" | `--font 3-d` or `--font bulbhead` |
| "big", "huge" | `--font big` |
| "readme header", "title" | `--border rounded` or `--border double`, `--align center` |
| "colorful", "gradient" | `--gradient "color1:color2"` (see below) |
| a brand/theme color | `--color` with the closest named color, or a hex code |

Default to `--font standard`, `--border none`, `--color default`,
`--align left` when the request is genuinely open-ended — those are the
tool's own defaults, so an unstyled first draft is a safe starting point to
iterate from.

This table is tuned for **terminal** output. If the result is headed to a
png/svg file instead, swap in the fonts/borders from the "Image output"
section below first — several terminal-safe choices above (`doom`, `block`,
`--border bold`) come out illegible or broken once rasterized as an image.

## Output format: pick based on where the result is going

- **Shown in chat / terminal** → default `terminal` output. Pass `--no-color`
  when the surrounding context won't render ANSI (e.g. embedding the result
  in a code block in a written response) — the CLI also auto-strips ANSI
  when stdout isn't a TTY, but don't rely on that when capturing output
  through a subprocess call.
- **Saved as an image (png/svg) for a README, slide, or social post** → see
  "Image output: font and border are not interchangeable with terminal
  choices" below *before* picking flags — the PNG renderer can't display
  every font/border that looks great in a terminal.
- **Dropped into a text file** → `--output txt`.
- **Consumed by another script/agent** → `--output json`; parse the `plain`
  field (always ANSI-free) rather than `styled`.

### Image output: font and border are not interchangeable with terminal choices

PNG rendering rasterizes glyphs through a bundled monospace font with
limited Unicode coverage, so a font or border that looks sharp in a
terminal can come out illegible or full of tofu (☐) boxes as an image.
Verified by rendering samples, not assumption:

- **Fonts that stay legible in PNG:** `ansi_shadow` (solid filled
  letters — the safest default for a "big bold logo"), `banner` (blocky,
  clean), `3-d` (textured but readable), `slant` (thin but readable).
- **Fonts to avoid for PNG/SVG:** `doom`, `block`, `standard`, `big`,
  `shadow`, `bulbhead`, `isometric1` — these render with overlapping or
  broken strokes once rasterized, even though they look fine in a
  terminal. Reach for `ansi_shadow` or `banner` instead when the user
  wants an image, even if `doom` felt like the right terminal answer for
  "bold."
- **Borders that stay legible in PNG:** `single`, `double`, `rounded`,
  `ascii`, `classic`, `shadow` (the drop-shadow style).
- **Borders to avoid for PNG/SVG:** `bold`, `dotted` — their box-drawing
  characters aren't in the bundled font and render as empty tofu boxes.

Always set `--bg-color` and `--color` to genuinely contrasting values
(default background is `black`) — a logo nobody can read is worse than an
unstyled one. If you generate a png/svg, glance at the file (or describe
what you'd expect) rather than assuming the terminal-tested combo carried
over cleanly.

For file outputs, use `--out-file <path>` when the user names a location or
filename convention; otherwise let the tool default to `./output.<ext>` in
the current directory and tell the user where it landed.

## Gradients and colors

- Named colors: `default red green blue cyan magenta yellow white`.
- Hex is also accepted: `--color "#FF6600"` or `--color "#F60"`.
- Gradients blend two colors across the text left-to-right:
  `--gradient "blue:cyan"` or `--gradient "#FF0000:#0000FF"`. Gradient and
  `--color` are mutually exclusive in effect — gradient wins if both are set.

## Examples

**Input:** "make an ascii banner for 'CHARM' with a double border in cyan"
```bash
charmascii "CHARM" --border double --color cyan
```

**Input:** "I want a big bold logo saying 'ACME' for my landing page, save as an image"
```bash
charmascii "ACME" --font ansi_shadow --color white --bg-color black --output png --out-file acme-logo.png
```
(`ansi_shadow` is used here instead of `doom`/`--border bold` specifically
*because* the output is an image — see the image-output section above.)

**Input:** "give me a rainbow-ish gradient title for my README, centered"
```bash
charmascii "My Project" --font slant --gradient "magenta:cyan" --align center --border rounded
```

**Input:** "I need this in JSON for my build script"
```bash
charmascii "v2.0" --output json
```

## Full reference

See `reference.md` for: complete flag table, all font names, border
styles, padding/width behavior, and the `charmascii mcp` server (for wiring
charmascii into an MCP-compatible agent directly instead of shelling out).
