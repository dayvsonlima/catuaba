package generator

import (
	"fmt"
	"os"
)

// IsInsideCatuabaProject checks if the current directory contains a valid Catuaba project.
func IsInsideCatuabaProject() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to determine current directory: %w", err)
	}

	if _, err := os.Stat(cwd + "/go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("no go.mod found — are you inside a Catuaba project?\n  Run 'catuaba new <name>' to create one")
	}

	if _, err := os.Stat(cwd + "/config/routes.go"); os.IsNotExist(err) {
		return fmt.Errorf("no config/routes.go found — this doesn't look like a Catuaba project.\n  Run 'catuaba new <name>' to create one")
	}

	return nil
}
