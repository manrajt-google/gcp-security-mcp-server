package tools

import (
	"context"
	"fmt"
	"os"

	scalibr "github.com/google/osv-scalibr"
	scalibrfs "github.com/google/osv-scalibr/fs"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func scalibrTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("scalibr",
			mcp.WithDescription("File system scanner used to extract software inventory data (e.g. installed language packages) and detect known vulnerabilities or generate SBOMs"),
			mcp.WithString("directory",
				mcp.Required(),
				mcp.Description("Directory to scan. This should be a valid path on the user's system."),
			),
		),
		Handler: scalibrHandler,
	}
}

func scalibrHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dir, err := request.RequireString("directory")
	if err != nil {
		return mcp.NewToolResultError("Missing directory: " + err.Error()), nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return mcp.NewToolResultError(fmt.Sprintf("Directory '%s' does not exist", dir)), nil
	}

	// plugins, err := pl.FromNames([]string{"os", "language", "secrets"})
	// if err != nil {
	// 	return mcp.NewToolResultError("Error loading plugins: " + err.Error()), nil
	// }

	cfg := &scalibr.ScanConfig{
		ScanRoots: scalibrfs.RealFSScanRoots(dir),
		// Plugins:   plugins,
	}

	results := scalibr.New().Scan(context.Background(), cfg)

	packages := []map[string]interface{}{}
	for _, i := range results.Inventory.Packages {
		packages = append(packages, map[string]interface{}{
			"name":    i.Name,
			"version": i.Version,
			"purl":    i.PURL,
		})
	}

	return mcp.NewToolResultStructuredOnly(map[string]interface{}{
		"packages": packages,
	}), nil
}
