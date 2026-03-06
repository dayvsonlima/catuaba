package code_editor

import (
	"regexp"
	"strings"
)

// InsertAttribute inserts a new attribute into an already-called method.
// It handles both single-line and multi-line calls, deduplicating existing attributes.
func InsertAttribute(code, methodName, newAttribute string) string {

	rgx := `(m?)\.` + methodName + `\(((.|\n)+?)\)`
	compiledRegex := regexp.MustCompile(rgx)

	submatch := compiledRegex.FindStringSubmatch(code)
	if submatch == nil {
		return code
	}

	existing := strings.TrimSpace(submatch[2])

	// Empty call — just insert
	if len(existing) == 0 {
		return compiledRegex.ReplaceAllString(code, `.`+methodName+`(`+newAttribute+`)`)
	}

	// Check if this attribute already exists (avoid duplicates)
	if strings.Contains(existing, newAttribute) {
		return code
	}

	// Strip trailing comma so we can add cleanly
	existing = strings.TrimRight(existing, ", \t\n")

	// Detect multi-line format
	if strings.Contains(existing, "\n") {
		// Multi-line: preserve indentation
		return compiledRegex.ReplaceAllString(code, `.`+methodName+`(`+"\n\t\t"+existing+",\n\t\t"+newAttribute+`,`+"\n\t"+`)`)
	}

	// Single-line
	return compiledRegex.ReplaceAllString(code, `.`+methodName+`(`+existing+`, `+newAttribute+`)`)
}
