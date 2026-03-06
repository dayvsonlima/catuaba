package mcp

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/dayvsonlima/catuaba/plugin"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func makeInstallPluginHandler(cache *Cache) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		source, err := req.RequireString("source")
		if err != nil {
			return mcpError("source is required"), nil
		}

		// Resolve plugin source
		pluginDir, cleanup, err := plugin.ResolveSource(source)
		if err != nil {
			return mcpError(fmt.Sprintf("failed to resolve plugin: %v", err)), nil
		}
		defer cleanup()

		// Load manifest
		m, err := plugin.LoadManifest(pluginDir + "/plugin.yaml")
		if err != nil {
			return mcpError(fmt.Sprintf("failed to load manifest: %v", err)), nil
		}

		// Validate
		if err := plugin.Validate(m); err != nil {
			return mcpError(fmt.Sprintf("invalid plugin: %v", err)), nil
		}

		// Install with no custom vars (MCP tools don't have interactive input)
		if err := plugin.Install(pluginDir, m, nil); err != nil {
			return mcpError(fmt.Sprintf("installation failed: %v", err)), nil
		}

		// Install Go dependencies
		for _, dep := range m.Dependencies {
			cmd := exec.Command("go", "get", dep)
			cmd.Dir = "."
			cmd.CombinedOutput() // best effort
		}

		// Track
		plugin.TrackInstall(m)

		cache.InvalidateAll()
		return mcpText(fmt.Sprintf("plugin %s v%s installed successfully", m.Name, m.Version)), nil
	}
}
