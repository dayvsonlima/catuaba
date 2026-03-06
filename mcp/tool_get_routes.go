package mcp

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dayvsonlima/catuaba/mcp/parser"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func makeGetRoutesHandler(projectDir string, cache *Cache) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		prefix := req.GetString("prefix", "")
		routesFile := filepath.Join(projectDir, "config/routes.go")

		routes, err := parser.ParseRoutes(routesFile)
		if err != nil {
			return mcpError(fmt.Sprintf("failed to parse routes: %v", err)), nil
		}

		if prefix != "" {
			var filtered []parser.RouteInfo
			for _, r := range routes {
				if strings.HasPrefix(r.Path, prefix) {
					filtered = append(filtered, r)
				}
			}
			return mcpJSON(filtered)
		}

		return mcpJSON(routes)
	}
}
