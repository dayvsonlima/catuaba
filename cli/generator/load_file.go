package generator

import "fmt"

func LoadFile(fileName string) string {
	output, err := Box.FindString(fileName)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return output
}
