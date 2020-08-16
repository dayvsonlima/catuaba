package code_editor

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func EditFile(filePath string, f func(content string) string) error {

	var (
		currentPath, _ = os.Getwd()
		fullPath       = currentPath + filePath
		code           = GetFileContent(filePath)
		result         = f(code)
	)

	err := ioutil.WriteFile(fullPath, []byte(result), 0644)

	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
		return err
	}

	return nil
}

func GetFileContent(filePath string) string {
	currentPath, _ := os.Getwd()
	content, err := ioutil.ReadFile(currentPath + filePath)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
