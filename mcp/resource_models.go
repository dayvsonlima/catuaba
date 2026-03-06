package mcp

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/dayvsonlima/catuaba/mcp/parser"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func handleModelsResource(projectDir string, cache *Cache) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcplib.ReadResourceRequest) ([]mcplib.ResourceContents, error) {
		modelsDir := filepath.Join(projectDir, "app/models")

		if cached, ok := cache.Get("models", modelsDir); ok {
			return cached.([]mcplib.ResourceContents), nil
		}

		models, err := parser.ParseModels(modelsDir)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(models)
		if err != nil {
			return nil, err
		}

		result := []mcplib.ResourceContents{
			mcplib.TextResourceContents{
				URI:      "catuaba://project/models",
				MIMEType: "application/json",
				Text:     string(data),
			},
		}
		cache.Set("models", result)
		return result, nil
	}
}
