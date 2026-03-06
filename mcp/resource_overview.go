package mcp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/dayvsonlima/catuaba/mcp/parser"
	"github.com/dayvsonlima/catuaba/plugin"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func handleOverviewResource(projectDir string, cache *Cache) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcplib.ReadResourceRequest) ([]mcplib.ResourceContents, error) {
		goModPath := filepath.Join(projectDir, "go.mod")

		if cached, ok := cache.Get("overview", goModPath); ok {
			return cached.([]mcplib.ResourceContents), nil
		}

		overview := buildOverview(projectDir, goModPath)
		data, err := json.Marshal(overview)
		if err != nil {
			return nil, err
		}

		result := []mcplib.ResourceContents{
			mcplib.TextResourceContents{
				URI:      "catuaba://project/overview",
				MIMEType: "application/json",
				Text:     string(data),
			},
		}
		cache.Set("overview", result)
		return result, nil
	}
}

func buildOverview(projectDir, goModPath string) parser.ProjectOverview {
	overview := parser.ProjectOverview{
		Module: "application",
		Go:     "1.22",
	}

	content, err := os.ReadFile(goModPath)
	if err == nil {
		for _, line := range strings.Split(string(content), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "module ") {
				overview.Module = strings.TrimPrefix(line, "module ")
			}
			if strings.HasPrefix(line, "go ") {
				overview.Go = strings.TrimPrefix(line, "go ")
			}
		}
	}

	envVars := parser.ParseEnv(projectDir)
	for _, v := range envVars {
		switch v {
		case "DB_DRIVER":
			overview.DB = readEnvValue(projectDir, "DB_DRIVER")
		case "APP_PORT":
			overview.Port = readEnvValue(projectDir, "APP_PORT")
		}
	}

	tracker, err := plugin.LoadTracker()
	if err == nil && len(tracker.Installed) > 0 {
		for name := range tracker.Installed {
			overview.Plugins = append(overview.Plugins, name)
		}
	}

	overview.Dirs = listDirs(projectDir)

	return overview
}

func readEnvValue(projectDir, key string) string {
	data, err := os.ReadFile(filepath.Join(projectDir, ".env"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, key+"=") {
			return strings.TrimPrefix(line, key+"=")
		}
	}
	return ""
}

func listDirs(projectDir string) []string {
	candidates := []string{
		"config",
		"database",
		"app/models",
		"app/controllers",
		"middleware",
	}
	var dirs []string

	controllersDir := filepath.Join(projectDir, "app/controllers")
	entries, err := os.ReadDir(controllersDir)
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				dirs = append(dirs, "app/controllers/"+e.Name())
			}
		}
	}

	for _, d := range candidates {
		if info, err := os.Stat(filepath.Join(projectDir, d)); err == nil && info.IsDir() {
			if d == "app/controllers" && len(dirs) > 0 {
				continue
			}
			dirs = append(dirs, d)
		}
	}

	return dirs
}
