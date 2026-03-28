// Command charmascii generates ASCII art from text.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/emmanuelgautier/charmascii"
	"github.com/emmanuelgautier/charmascii/internal/output"
)

// Populated by GoReleaser via -ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var (
		font        string
		borderStyle string
		textColor   string
		borderColor string
		align       string
		padding     int
		vPadding    int
		outputFmt   string
		outFile     string
		width       int
		gradient    string
		bgColor     string
		noColor     bool
		textShadow  bool
		listFonts   bool
		showVersion bool
	)

	cmd := &cobra.Command{
		Use:   "charmascii [text]",
		Short: "Generate ASCII art from text",
		Long: `charmascii converts text to ASCII art using FIGlet fonts, with
optional borders, colors, gradients, and multiple output formats.`,
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Printf("charmascii %s (%s) built %s\n", version, commit, date)
				return nil
			}

			if listFonts {
				for _, f := range charmascii.ListFonts() {
					fmt.Println(f)
				}
				return nil
			}

			if len(args) == 0 {
				return cmd.Help()
			}

			text := args[0]
			opts := charmascii.Options{
				Font:        font,
				Border:      borderStyle,
				Color:       textColor,
				BorderColor: borderColor,
				Align:       align,
				Padding:     padding,
				VPadding:    vPadding,
				Width:       width,
				Gradient:    gradient,
				BgColor:     bgColor,
				NoColor:     noColor || !output.IsTerminal(os.Stdout),
				TextShadow:  textShadow,
			}

			result, err := charmascii.Generate(text, opts)
			if err != nil {
				return err
			}

			// Determine output file name.
			dest := outFile
			if dest == "" && outputFmt != "terminal" {
				dest = "output." + defaultExt(outputFmt)
			}

			meta := output.Metadata{
				Version: version,
				Command: buildCommand(cmd, args[0], font, borderStyle, textColor, borderColor, align,
					padding, vPadding, width, outputFmt, outFile, gradient, bgColor, noColor, textShadow),
				URL: "https://github.com/emmanuelgautier/charmascii",
			}

			switch outputFmt {
			case "terminal", "":
				return output.WriteTerminal(os.Stdout, strings.Split(result.Styled, "\n"), opts.NoColor)
			case "txt":
				return output.WriteTXT(dest, result.Lines, meta)
			case "png":
				return output.WritePNG(dest, result.Lines, bgColor, textColor, meta)
			case "svg":
				return output.WriteSVG(dest, result.Lines, bgColor, textColor, meta)
			default:
				return fmt.Errorf("unknown output format %q; valid choices: terminal, txt, png, svg", outputFmt)
			}
		},
	}

	f := cmd.Flags()
	f.StringVar(&font, "font", "standard",
		"FIGlet font (standard|big|doom|isometric1|slant|block|3-d|shadow|banner|bulbhead)")
	f.StringVar(&borderStyle, "border", "none",
		"Border style (none|single|double|rounded|bold|ascii)")
	f.StringVar(&textColor, "color", "default",
		"Text color (default|red|green|blue|cyan|magenta|yellow|white|#RRGGBB|#RGB)")
	f.StringVar(&borderColor, "border-color", "default",
		"Border color (same choices as --color)")
	f.StringVar(&align, "align", "left",
		"Text alignment (left|center|right)")
	f.IntVar(&padding, "padding", 1,
		"Inner horizontal padding inside the border box")
	f.IntVar(&vPadding, "v-padding", -1,
		"Inner vertical padding (blank lines) inside the border box (default: same as --padding)")
	f.StringVar(&outputFmt, "output", "terminal",
		"Output format (terminal|txt|png|svg)")
	f.StringVar(&outFile, "out-file", "",
		"Output file path (default: ./output.<ext>)")
	f.IntVar(&width, "width", 0,
		"Max width in characters (default: terminal width)")
	f.StringVar(&gradient, "gradient", "",
		`Two-color gradient, e.g. "red:blue", "green:cyan", or "#FF0000:#0000FF"`)
	f.StringVar(&bgColor, "bg-color", "black",
		"Background color for PNG/SVG output")
	f.BoolVar(&noColor, "no-color", false,
		"Strip all ANSI codes (auto-enabled when stdout is not a TTY)")
	f.BoolVar(&textShadow, "text-shadow", false,
		"Add a drop shadow (░) behind the ASCII-art letters")
	f.BoolVar(&listFonts, "list-fonts", false,
		"Print all available fonts and exit")
	f.BoolVar(&showVersion, "version", false,
		"Print version information and exit")

	return cmd
}

// buildCommand reconstructs the CLI invocation from explicitly-set flags.
func buildCommand(cmd *cobra.Command, text, font, borderStyle, textColor, borderColor, align string,
	padding, vPadding, width int, outputFmt, outFile, gradient, bgColor string, noColor, textShadow bool,
) string {
	parts := []string{"charmascii", fmt.Sprintf("%q", text)}
	f := cmd.Flags()
	if f.Changed("font") {
		parts = append(parts, "--font", font)
	}
	if f.Changed("border") {
		parts = append(parts, "--border", borderStyle)
	}
	if f.Changed("color") {
		parts = append(parts, "--color", textColor)
	}
	if f.Changed("border-color") {
		parts = append(parts, "--border-color", borderColor)
	}
	if f.Changed("align") {
		parts = append(parts, "--align", align)
	}
	if f.Changed("padding") {
		parts = append(parts, "--padding", fmt.Sprint(padding))
	}
	if f.Changed("v-padding") {
		parts = append(parts, "--v-padding", fmt.Sprint(vPadding))
	}
	if f.Changed("output") {
		parts = append(parts, "--output", outputFmt)
	}
	if f.Changed("out-file") {
		parts = append(parts, "--out-file", outFile)
	}
	if f.Changed("width") {
		parts = append(parts, "--width", fmt.Sprint(width))
	}
	if f.Changed("gradient") {
		parts = append(parts, "--gradient", gradient)
	}
	if f.Changed("bg-color") {
		parts = append(parts, "--bg-color", bgColor)
	}
	if f.Changed("no-color") && noColor {
		parts = append(parts, "--no-color")
	}
	if f.Changed("text-shadow") && textShadow {
		parts = append(parts, "--text-shadow")
	}
	return strings.Join(parts, " ")
}

func defaultExt(format string) string {
	switch format {
	case "txt":
		return "txt"
	case "png":
		return "png"
	case "svg":
		return "svg"
	default:
		return "out"
	}
}
