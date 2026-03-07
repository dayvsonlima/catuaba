package mcp

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/dayvsonlima/catuaba/mcp/parser"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func handleComponentsResource(projectDir string, cache *Cache) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcplib.ReadResourceRequest) ([]mcplib.ResourceContents, error) {
		componentsDir := filepath.Join(projectDir, "app/views/components")

		if cached, ok := cache.Get("components", componentsDir); ok {
			return cached.([]mcplib.ResourceContents), nil
		}

		components, err := parser.ParseComponents(componentsDir)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(components)
		if err != nil {
			return nil, err
		}

		result := []mcplib.ResourceContents{
			mcplib.TextResourceContents{
				URI:      "catuaba://project/components",
				MIMEType: "application/json",
				Text:     string(data),
			},
		}
		cache.Set("components", result)
		return result, nil
	}
}
