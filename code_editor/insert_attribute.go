package code_editor

import (
	"regexp"
	"strings"
)

// InsertAttribute
// insert a new attribute in a already called method
func InsertAttribute(code, methodName, newAttribute string) string {

	rgx := `(m?)\.` + methodName + `\(((.|\n|)+)\)`
	compiledRegex := regexp.MustCompile(rgx)

	submatch := compiledRegex.FindStringSubmatch(code)
	hasAttributes := len(strings.Trim(submatch[2], " "))

	if hasAttributes == 0 {
		return compiledRegex.ReplaceAllString(code, `.`+methodName+`(`+newAttribute+`)`)
	}

	return compiledRegex.ReplaceAllString(code, `.`+methodName+`($2, `+newAttribute+`)`)
}
