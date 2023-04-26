package generator

import (
	"fmt"
	"io/ioutil"
	"os"
)

func GenerateFile(template string, data interface{}, filePath string) error {
	currentPath, _ := os.Getwd()
	fullPath := currentPath + "/" + filePath

	content := Render(template, data)
	err := ioutil.WriteFile(fullPath, []byte(content), 0644)

	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
		return err
	}

	fmt.Println(filePath)

	return nil
}

func GenerateFromContent(content string, data interface{}, filePath string) error {
	currentPath, _ := os.Getwd()
	fullPath := currentPath + "/" + filePath

	result := RenderFromContent(content, data)
	err := ioutil.WriteFile(fullPath, []byte(result), 0644)

	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
		return err
	}

	fmt.Println(filePath)

	return nil
}
