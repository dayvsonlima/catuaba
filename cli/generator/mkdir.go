package generator

import (
	"fmt"
	"os"
)

func Mkdir(path string) {
	currentPath, _ := os.Getwd()
	fullPath := currentPath + "/" + path

	fmt.Println(fullPath)

	_ = os.Mkdir(fullPath, 0755)
}
