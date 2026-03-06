package code_editor

import (
	"fmt"
	"os"
)

func EditFile(filePath string, f func(content string) string) error {
	currentPath, _ := os.Getwd()
	fullPath := currentPath + "/" + filePath

	code, err := GetFileContent(filePath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	result := f(code)

	if err := os.WriteFile(fullPath, []byte(result), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", filePath, err)
	}

	return nil
}

func GetFileContent(filePath string) (string, error) {
	currentPath, _ := os.Getwd()
	content, err := os.ReadFile(currentPath + "/" + filePath)
	if err != nil {
		return "", fmt.Errorf("unable to read file %s: %w", filePath, err)
	}
	return string(content), nil
}
