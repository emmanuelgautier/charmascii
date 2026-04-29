package main

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	"github.com/emmanuelgautier/charmascii/internal/mcpserver"
)

func newMCPCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "Start an MCP server over stdio",
		Long:  "Starts a JSON-RPC 2.0 MCP server on stdin/stdout exposing the generate_ascii tool.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := mcpserver.New(version)
			return server.ServeStdio(s)
		},
	}
}
