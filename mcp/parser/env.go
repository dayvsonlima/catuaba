package parser

import (
	"os"
	"path/filepath"
	"strings"
)

// ParseEnv reads .env (or .env.example) and returns the list of variable names (no values).
func ParseEnv(projectDir string) []string {
	// Try .env first, then .env.example
	candidates := []string{".env", ".env.example"}
	for _, name := range candidates {
		vars := parseEnvFile(filepath.Join(projectDir, name))
		if len(vars) > 0 {
			return vars
		}
	}
	return nil
}

func parseEnvFile(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var vars []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) >= 1 {
			key := strings.TrimSpace(parts[0])
			if key != "" {
				vars = append(vars, key)
			}
		}
	}
	return vars
}
