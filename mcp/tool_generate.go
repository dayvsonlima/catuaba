package mcp

import (
	"context"
	"strings"

	"github.com/dayvsonlima/catuaba/cli/api"
	"github.com/dayvsonlima/catuaba/cli/model"
	"github.com/dayvsonlima/catuaba/cli/scaffold"
	"github.com/dayvsonlima/catuaba/generator"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func makeGenerateScaffoldHandler(cache *Cache) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		name, err := req.RequireString("name")
		if err != nil {
			return mcpError("name is required"), nil
		}
		fieldsStr, err := req.RequireString("fields")
		if err != nil {
			return mcpError("fields is required"), nil
		}

		fields := strings.Split(fieldsStr, ",")
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}

		data := model.ModelBuilder{
			Name:   generator.Camelize(name),
			Params: fields,
		}

		if err := model.BuildModel(data); err != nil {
			return mcpError("model generation failed: " + err.Error()), nil
		}
		if err := scaffold.BuildHandlers(data); err != nil {
			return mcpError("handler generation failed: " + err.Error()), nil
		}
		if err := scaffold.BuildViews(data); err != nil {
			return mcpError("view generation failed: " + err.Error()), nil
		}
		if err := scaffold.BuildRoutes(data); err != nil {
			return mcpError("routes generation failed: " + err.Error()), nil
		}
		if err := api.BuildControllers(data); err != nil {
			return mcpError("api controller generation failed: " + err.Error()), nil
		}
		if err := api.BuildRoutes(data); err != nil {
			return mcpError("api routes generation failed: " + err.Error()), nil
		}

		cache.InvalidateAll()
		return mcpText("scaffold generated successfully for " + data.Name), nil
	}
}

func makeGenerateModelHandler(cache *Cache) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		name, err := req.RequireString("name")
		if err != nil {
			return mcpError("name is required"), nil
		}
		fieldsStr, err := req.RequireString("fields")
		if err != nil {
			return mcpError("fields is required"), nil
		}

		fields := strings.Split(fieldsStr, ",")
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}

		data := model.ModelBuilder{
			Name:   generator.Camelize(name),
			Params: fields,
		}

		if err := model.BuildModel(data); err != nil {
			return mcpError("model generation failed: " + err.Error()), nil
		}

		cache.InvalidateAll()
		return mcpText("model generated successfully: " + data.Name), nil
	}
}
