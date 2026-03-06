package code_editor_test

import (
	"testing"

	"github.com/dayvsonlima/catuaba/code_editor"
	"github.com/stretchr/testify/assert"
)

func TestAddImport(t *testing.T) {

	t.Run("adds a new import to existing imports", func(t *testing.T) {
		code := `package fake_package

import (
	"fmt"
)
`
		newPkg := "github.com/urfave/cli/v2"
		actual := code_editor.AddImport(code, newPkg)

		assert.Contains(t, actual, newPkg)
		assert.Contains(t, actual, "fmt")
	})

	t.Run("does not duplicate existing import", func(t *testing.T) {
		code := `package fake_package

import (
	"fmt"
)
`
		actual := code_editor.AddImport(code, "fmt")
		assert.Contains(t, actual, "fmt")
	})

	t.Run("handles multiple existing imports", func(t *testing.T) {
		code := `package fake_package

import (
	"fmt"
	"net/http"
)
`
		actual := code_editor.AddImport(code, "os")
		assert.Contains(t, actual, "fmt")
		assert.Contains(t, actual, "net/http")
		assert.Contains(t, actual, "os")
	})
}

func TestAddImportIfNotExist(t *testing.T) {

	code := `package fake_package

import (
	"fmt"
)
`

	t.Run("adds import when not present", func(t *testing.T) {
		newPkg := "github.com/urfave/cli/v2"
		actual := code_editor.AddImportIfNotExist(code, newPkg)

		assert.Contains(t, actual, newPkg)
		assert.Contains(t, actual, "fmt")
	})

	t.Run("should not add if already exists", func(t *testing.T) {
		actual := code_editor.AddImportIfNotExist(code, "fmt")
		assert.Equal(t, actual, code)
	})
}
