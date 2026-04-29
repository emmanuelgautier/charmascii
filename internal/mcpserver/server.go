package mcpserver

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	charmascii "github.com/emmanuelgautier/charmascii"
)

// New builds an MCPServer that exposes the generate_ascii tool.
func New(version string) *server.MCPServer {
	s := server.NewMCPServer("charmascii", version)

	tool := mcp.NewTool("generate_ascii",
		mcp.WithDescription("Generate ASCII art from text using FIGlet fonts"),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("Text to convert to ASCII art"),
		),
		mcp.WithString("font",
			mcp.DefaultString("standard"),
			mcp.Description("FIGlet font: standard|big|doom|isometric1|slant|block|3-d|shadow|banner|bulbhead|ansi_shadow"),
		),
		mcp.WithString("border",
			mcp.DefaultString("none"),
			mcp.Description("Border style: none|single|double|rounded|bold|ascii|classic|dotted|shadow"),
		),
		mcp.WithString("color",
			mcp.DefaultString("default"),
			mcp.Description("Text color: default|red|green|blue|cyan|magenta|yellow|white|#RRGGBB"),
		),
		mcp.WithString("align",
			mcp.DefaultString("left"),
			mcp.Description("Text alignment: left|center|right"),
		),
		mcp.WithNumber("padding",
			mcp.DefaultNumber(1),
			mcp.Description("Inner horizontal padding inside the border box"),
		),
		mcp.WithNumber("width",
			mcp.DefaultNumber(0),
			mcp.Description("Max width in characters; 0 means no limit"),
		),
		mcp.WithString("gradient",
			mcp.DefaultString(""),
			mcp.Description(`Two-color gradient e.g. "red:blue" or "#FF0000:#0000FF"`),
		),
		mcp.WithBoolean("text_shadow",
			mcp.DefaultBool(false),
			mcp.Description("Add a drop shadow (░) behind the ASCII art letters"),
		),
	)

	s.AddTool(tool, generateHandler)
	return s
}

func generateHandler(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text := req.GetString("text", "")
	if text == "" {
		return mcp.NewToolResultError("text parameter is required"), nil
	}

	opts := charmascii.Options{
		Font:       req.GetString("font", "standard"),
		Border:     req.GetString("border", "none"),
		Color:      req.GetString("color", "default"),
		Align:      req.GetString("align", "left"),
		Padding:    req.GetInt("padding", 1),
		Width:      req.GetInt("width", 0),
		Gradient:   req.GetString("gradient", ""),
		TextShadow: req.GetBool("text_shadow", false),
		NoColor:    true, // LLM clients cannot render ANSI escape codes
	}

	result, err := charmascii.Generate(text, opts)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(strings.Join(result.Lines, "\n")), nil
}
