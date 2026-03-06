package mcp

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/dayvsonlima/catuaba/mcp/parser"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func handleRoutesResource(projectDir string, cache *Cache) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcplib.ReadResourceRequest) ([]mcplib.ResourceContents, error) {
		routesFile := filepath.Join(projectDir, "config/routes.go")

		if cached, ok := cache.Get("routes", routesFile); ok {
			return cached.([]mcplib.ResourceContents), nil
		}

		routes, err := parser.ParseRoutes(routesFile)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(routes)
		if err != nil {
			return nil, err
		}

		result := []mcplib.ResourceContents{
			mcplib.TextResourceContents{
				URI:      "catuaba://project/routes",
				MIMEType: "application/json",
				Text:     string(data),
			},
		}
		cache.Set("routes", result)
		return result, nil
	}
}
