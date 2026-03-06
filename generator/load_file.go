package generator

import "fmt"

func LoadFile(fileName string) string {
	content, err := templateFS.ReadFile("templates/" + fileName)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(content)
}
