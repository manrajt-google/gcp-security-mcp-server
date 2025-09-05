package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func greetTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("greet",
			mcp.WithDescription("Greets the user"),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Name of user"),
			),
		),
		Handler: greetHandler,
	}
}

func greetHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError("Missing name: " + err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
