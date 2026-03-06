package generator

import (
	"fmt"
	"os"

	"github.com/dayvsonlima/catuaba/cli/output"
)

func GenerateFile(tmpl string, data interface{}, filePath string) error {
	currentPath, _ := os.Getwd()
	fullPath := currentPath + "/" + filePath

	content, err := Render(tmpl, data)
	if err != nil {
		return fmt.Errorf("rendering %s: %w", tmpl, err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		output.Error("Unable to write file: %v", err)
		return err
	}

	output.Create(filePath)
	return nil
}

func GenerateFromContent(content string, data interface{}, filePath string) error {
	currentPath, _ := os.Getwd()
	fullPath := currentPath + "/" + filePath

	result, err := RenderFromContent(content, data)
	if err != nil {
		return fmt.Errorf("rendering content for %s: %w", filePath, err)
	}

	if err := os.WriteFile(fullPath, []byte(result), 0644); err != nil {
		output.Error("Unable to write file: %v", err)
		return err
	}

	output.Create(filePath)
	return nil
}
