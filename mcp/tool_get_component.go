package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dayvsonlima/catuaba/mcp/parser"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func makeGetComponentHandler(projectDir string, cache *Cache) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		name := req.GetString("name", "")
		componentsDir := filepath.Join(projectDir, "app/views/components")

		result, err := parser.ParseComponents(componentsDir)
		if err != nil {
			return mcpError(fmt.Sprintf("failed to parse components: %v", err)), nil
		}

		if len(result.Components) == 0 {
			return mcpError("no components found in app/views/components/"), nil
		}

		// No name: return compact list
		if name == "" {
			return mcpJSON(result)
		}

		// Find matching component and return its source file
		for _, c := range result.Components {
			if strings.EqualFold(c.Name, name) {
				filePath := filepath.Join(componentsDir, c.File)
				src, err := os.ReadFile(filePath)
				if err != nil {
					return mcpError(fmt.Sprintf("failed to read %s: %v", c.File, err)), nil
				}
				return mcpText(string(src)), nil
			}
		}

		return mcpError(fmt.Sprintf("component %q not found", name)), nil
	}
}
