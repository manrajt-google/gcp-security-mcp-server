package tools

import (
	"context"

	// Todo: Figure out if we should change mark3labs mcp package to the official mcp go sdk package
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Add adds the tools to the mcp server.
func Add(_ context.Context, s *server.MCPServer) {
	s.AddTools(
		greetTool(),
		scalibrTool(),
	)
}

func result(s string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(s),
		},
	}
}
