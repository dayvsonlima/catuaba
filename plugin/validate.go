package plugin

import (
	"fmt"
	"os"
)

// Validate checks the manifest for required fields and potential conflicts
// against the current project.
func Validate(m *Manifest) error {
	if m.Name == "" {
		return fmt.Errorf("plugin manifest missing required field: name")
	}
	if m.Version == "" {
		return fmt.Errorf("plugin manifest missing required field: version")
	}

	for i, f := range m.Files {
		if f.Template == "" {
			return fmt.Errorf("plugin file entry %d missing template path", i)
		}
		if f.Output == "" {
			return fmt.Errorf("plugin file entry %d missing output path", i)
		}
	}

	for i, inj := range m.Injections {
		if inj.File == "" {
			return fmt.Errorf("injection %d missing target file", i)
		}
		if inj.Action == "" {
			return fmt.Errorf("injection %d missing action", i)
		}
		switch inj.Action {
		case "add_import":
			if inj.Import == "" {
				return fmt.Errorf("injection %d (add_import) missing import path", i)
			}
		case "insert_attribute":
			if inj.Method == "" || inj.Attribute == "" {
				return fmt.Errorf("injection %d (insert_attribute) missing method or attribute", i)
			}
		case "insert_line_after":
			if inj.After == "" || inj.Content == "" {
				return fmt.Errorf("injection %d (insert_line_after) missing after or content", i)
			}
		case "insert_before_closing_brace":
			if inj.Content == "" {
				return fmt.Errorf("injection %d (insert_before_closing_brace) missing content", i)
			}
		case "append":
			if inj.Content == "" {
				return fmt.Errorf("injection %d (append) missing content", i)
			}
		default:
			return fmt.Errorf("injection %d has unknown action %q", i, inj.Action)
		}
	}

	return nil
}

// CheckConflicts checks if generated files already exist in the project.
// Returns a list of conflicting file paths.
func CheckConflicts(m *Manifest) []string {
	cwd, _ := os.Getwd()
	var conflicts []string
	for _, f := range m.Files {
		path := cwd + "/" + f.Output
		if _, err := os.Stat(path); err == nil {
			conflicts = append(conflicts, f.Output)
		}
	}
	return conflicts
}
