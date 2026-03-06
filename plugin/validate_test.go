package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate_Valid(t *testing.T) {
	m := &Manifest{
		Name:    "test",
		Version: "1.0.0",
		Files: []File{
			{Template: "t.go.tmpl", Output: "t.go"},
		},
		Injections: []Injection{
			{File: "routes.go", Action: "add_import", Import: "pkg/test"},
			{File: "routes.go", Action: "insert_attribute", Method: "AutoMigrate", Attribute: "&Test{}"},
			{File: "routes.go", Action: "insert_line_after", After: "import (", Content: `"new/pkg"`},
			{File: "routes.go", Action: "insert_before_closing_brace", Content: "r.GET(\"/\")"},
			{File: "routes.go", Action: "append", Content: "// end"},
		},
	}
	assert.NoError(t, Validate(m))
}

func TestValidate_MissingName(t *testing.T) {
	m := &Manifest{Version: "1.0.0"}
	assert.EqualError(t, Validate(m), "plugin manifest missing required field: name")
}

func TestValidate_MissingVersion(t *testing.T) {
	m := &Manifest{Name: "test"}
	assert.EqualError(t, Validate(m), "plugin manifest missing required field: version")
}

func TestValidate_UnknownAction(t *testing.T) {
	m := &Manifest{
		Name:    "test",
		Version: "1.0.0",
		Injections: []Injection{
			{File: "f.go", Action: "unknown_action"},
		},
	}
	assert.Contains(t, Validate(m).Error(), "unknown action")
}

func TestValidate_AddImportMissingPath(t *testing.T) {
	m := &Manifest{
		Name:    "test",
		Version: "1.0.0",
		Injections: []Injection{
			{File: "f.go", Action: "add_import"},
		},
	}
	assert.Contains(t, Validate(m).Error(), "missing import path")
}

func TestValidate_InsertAttributeMissingMethod(t *testing.T) {
	m := &Manifest{
		Name:    "test",
		Version: "1.0.0",
		Injections: []Injection{
			{File: "f.go", Action: "insert_attribute", Attribute: "&M{}"},
		},
	}
	assert.Contains(t, Validate(m).Error(), "missing method or attribute")
}
