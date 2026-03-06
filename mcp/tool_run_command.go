package mcp

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func makeRunCommandHandler(cache *Cache) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		command, err := req.RequireString("command")
		if err != nil {
			return mcpError("command is required"), nil
		}
		argsStr, err := req.RequireString("args")
		if err != nil {
			return mcpError("args is required"), nil
		}

		// Validate command is a known generate subcommand
		allowed := map[string]bool{
			"controller": true,
			"model":      true,
			"scaffold":   true,
			"migration":  true,
			"middleware":  true,
			"service":    true,
			"seed":       true,
		}
		if !allowed[command] {
			return mcpError(fmt.Sprintf("unknown command %q — allowed: controller, model, scaffold, migration, middleware, service, seed", command)), nil
		}

		args := strings.Split(argsStr, ",")
		for i := range args {
			args[i] = strings.TrimSpace(args[i])
		}

		// Build command: catuaba generate <command> <args...>
		cmdArgs := append([]string{"generate", command}, args...)
		cmd := exec.Command("catuaba", cmdArgs...)

		output, runErr := cmd.CombinedOutput()
		if runErr != nil {
			return mcpError(fmt.Sprintf("command failed: %v\n%s", runErr, string(output))), nil
		}

		cache.InvalidateAll()
		return mcpText(string(output)), nil
	}
}
