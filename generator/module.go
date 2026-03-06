package generator

import (
	"os"
	"strings"
)

// ModuleName reads the module name from the go.mod file in the current directory
func ModuleName() string {
	currentPath, _ := os.Getwd()
	content, err := os.ReadFile(currentPath + "/go.mod")
	if err != nil {
		return "application"
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}

	return "application"
}
