package plugin

import (
	"strings"
)

// InsertLineAfter finds the first line containing `after` and inserts `content`
// on the next line.
func InsertLineAfter(code, after, content string) string {
	lines := strings.Split(code, "\n")
	var result []string

	found := false
	for _, line := range lines {
		result = append(result, line)
		if !found && strings.Contains(line, after) {
			for _, cl := range strings.Split(content, "\n") {
				result = append(result, cl)
			}
			found = true
		}
	}

	return strings.Join(result, "\n")
}

// InsertBeforeClosingBrace inserts content before the last closing brace `}` in the code.
// This is the same pattern used by scaffold route injection.
func InsertBeforeClosingBrace(code, content string) string {
	return strings.ReplaceAll(code, "\n}", "\n"+content+"\n}")
}

// AppendToFile appends content to the end of the code, ensuring a newline separator.
func AppendToFile(code, content string) string {
	if !strings.HasSuffix(code, "\n") {
		code += "\n"
	}
	return code + content + "\n"
}
