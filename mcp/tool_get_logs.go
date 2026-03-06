package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dayvsonlima/catuaba/mcp/parser"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func makeGetLogsHandler(projectDir string) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		maxLines := req.GetInt("lines", 50)
		level := shortLevel(req.GetString("level", ""))
		pathFilter := req.GetString("path", "")

		logFile := filepath.Join(projectDir, "tmp/app.log")
		f, err := os.Open(logFile)
		if err != nil {
			return mcpError(fmt.Sprintf("no log file found at tmp/app.log: %v", err)), nil
		}
		defer f.Close()

		var allLines []string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			allLines = append(allLines, scanner.Text())
		}

		var entries []parser.LogEntry
		for i := len(allLines) - 1; i >= 0 && len(entries) < maxLines; i-- {
			entry := parseLine(allLines[i])
			if entry == nil {
				continue
			}
			if level != "" && !strings.EqualFold(entry.Level, level) {
				continue
			}
			if pathFilter != "" && !strings.HasPrefix(entry.Path, pathFilter) {
				continue
			}
			entries = append(entries, *entry)
		}

		// Reverse to chronological order
		for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
			entries[i], entries[j] = entries[j], entries[i]
		}

		return mcpJSON(entries)
	}
}

func parseLine(line string) *parser.LogEntry {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	var raw map[string]any
	if err := json.Unmarshal([]byte(line), &raw); err == nil {
		return parseJSONLog(raw)
	}

	return &parser.LogEntry{Message: line}
}

func parseJSONLog(raw map[string]any) *parser.LogEntry {
	entry := &parser.LogEntry{}

	if v, ok := raw["time"].(string); ok {
		if len(v) >= 19 {
			entry.Time = v[11:19]
		} else {
			entry.Time = v
		}
	}
	if v, ok := raw["level"].(string); ok {
		entry.Level = shortLevel(v)
	}
	if v, ok := raw["status"].(float64); ok {
		entry.Status = int(v)
	}
	if v, ok := raw["method"].(string); ok {
		entry.Method = v
	}
	if v, ok := raw["path"].(string); ok {
		entry.Path = v
	}
	if v, ok := raw["latency"].(string); ok {
		entry.Dur = v
	}
	if v, ok := raw["duration"].(string); ok {
		entry.Dur = v
	}
	if v, ok := raw["error"].(string); ok {
		entry.Error = v
	}
	if v, ok := raw["msg"].(string); ok {
		entry.Message = v
	}

	return entry
}

func shortLevel(level string) string {
	switch strings.ToUpper(level) {
	case "ERROR", "ERR":
		return "ERR"
	case "WARN", "WARNING":
		return "WRN"
	case "INFO":
		return "INF"
	case "DEBUG":
		return "DBG"
	default:
		if len(level) >= 3 {
			return strings.ToUpper(level[:3])
		}
		return strings.ToUpper(level)
	}
}
