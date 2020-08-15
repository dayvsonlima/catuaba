package code_editor

import (
	"regexp"
)

func InsertAttribute(code, methodName, newAttribute string) string {

	rgx := `(m?)\.` + methodName + `\((.+)\)`
	compiledRegex := regexp.MustCompile(rgx)
	output := compiledRegex.ReplaceAllString(code, `.`+methodName+`($2, `+newAttribute+`)`)

	return output
}
