package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/manrajt/gcp-security-mcp/internal/tools"
	"github.com/mark3labs/mcp-go/server"
)

// StartServer initializes and starts the MCP server.
func StartServer() {
	s := server.NewMCPServer(
		"GCP Security MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)
	slog.Info("Adding tools and resources to the server.")
	ctx := context.Background()
	tools.Add(ctx, s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
