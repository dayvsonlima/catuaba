package mcp

import (
	"context"
	"encoding/json"

	"github.com/dayvsonlima/catuaba/mcp/parser"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func handleEnvResource(projectDir string, cache *Cache) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcplib.ReadResourceRequest) ([]mcplib.ResourceContents, error) {
		if cached, ok := cache.Get("env"); ok {
			return cached.([]mcplib.ResourceContents), nil
		}

		envVars := parser.ParseEnv(projectDir)

		data, err := json.Marshal(envVars)
		if err != nil {
			return nil, err
		}

		result := []mcplib.ResourceContents{
			mcplib.TextResourceContents{
				URI:      "catuaba://project/env",
				MIMEType: "application/json",
				Text:     string(data),
			},
		}
		cache.Set("env", result)
		return result, nil
	}
}
