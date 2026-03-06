package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dayvsonlima/catuaba/mcp/parser"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func makeGetModelHandler(projectDir string, cache *Cache) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		name := req.GetString("name", "")
		modelsDir := filepath.Join(projectDir, "app/models")

		models, err := parser.ParseModels(modelsDir)
		if err != nil {
			return mcpError(fmt.Sprintf("failed to parse models: %v", err)), nil
		}

		if name != "" {
			for _, m := range models {
				if strings.EqualFold(m.Name, name) {
					return mcpJSON(m)
				}
			}
			return mcpError(fmt.Sprintf("model %q not found", name)), nil
		}

		return mcpJSON(models)
	}
}

func mcpJSON(v any) (*mcplib.CallToolResult, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return &mcplib.CallToolResult{
		Content: []mcplib.Content{
			mcplib.TextContent{
				Type: "text",
				Text: string(data),
			},
		},
	}, nil
}

func mcpText(text string) *mcplib.CallToolResult {
	return &mcplib.CallToolResult{
		Content: []mcplib.Content{
			mcplib.TextContent{
				Type: "text",
				Text: text,
			},
		},
	}
}

func mcpError(msg string) *mcplib.CallToolResult {
	return &mcplib.CallToolResult{
		Content: []mcplib.Content{
			mcplib.TextContent{
				Type: "text",
				Text: msg,
			},
		},
		IsError: true,
	}
}
