package mcp

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/dayvsonlima/catuaba/mcp/parser"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func handleControllersResource(projectDir string, cache *Cache) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcplib.ReadResourceRequest) ([]mcplib.ResourceContents, error) {
		controllersDir := filepath.Join(projectDir, "app/controllers")

		if cached, ok := cache.Get("controllers", controllersDir); ok {
			return cached.([]mcplib.ResourceContents), nil
		}

		controllers, err := parser.ParseControllers(controllersDir)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(controllers)
		if err != nil {
			return nil, err
		}

		result := []mcplib.ResourceContents{
			mcplib.TextResourceContents{
				URI:      "catuaba://project/controllers",
				MIMEType: "application/json",
				Text:     string(data),
			},
		}
		cache.Set("controllers", result)
		return result, nil
	}
}
