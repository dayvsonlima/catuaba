package mcp

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/dayvsonlima/catuaba/mcp/parser"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func handleMiddlewareResource(projectDir string, cache *Cache) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcplib.ReadResourceRequest) ([]mcplib.ResourceContents, error) {
		middlewareDir := filepath.Join(projectDir, "middleware")

		if cached, ok := cache.Get("middleware", middlewareDir); ok {
			return cached.([]mcplib.ResourceContents), nil
		}

		middlewares, err := parser.ParseMiddleware(middlewareDir)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(middlewares)
		if err != nil {
			return nil, err
		}

		result := []mcplib.ResourceContents{
			mcplib.TextResourceContents{
				URI:      "catuaba://project/middleware",
				MIMEType: "application/json",
				Text:     string(data),
			},
		}
		cache.Set("middleware", result)
		return result, nil
	}
}
