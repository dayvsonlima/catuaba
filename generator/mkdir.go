package generator

import (
	"fmt"
	"os"

	"github.com/dayvsonlima/catuaba/cli/output"
)

func Mkdir(path string) error {
	currentPath, _ := os.Getwd()
	fullPath := currentPath + "/" + path

	output.Mkdir(path)

	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", path, err)
	}
	return nil
}
