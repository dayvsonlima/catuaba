package code_editor_test

import (
	"testing"

	"github.com/dayvsonlima/catuaba/code_editor"
	"github.com/stretchr/testify/assert"
)

func TestAddImport(t *testing.T) {

	code := `package fake_package

	import (
		"fmt"
	)`

	newPkg := "github.com/urfave/cli/v2"

	actual := code_editor.AddImport(code, newPkg)

	assert.Contains(t, actual, newPkg)
	assert.Contains(t, actual, "fmt")
}

func TestAddImportIfNotExist(t *testing.T) {

	code := `package fake_package

  import (
    "fmt"
  )`

	newPkg := "github.com/urfave/cli/v2"

	actual := code_editor.AddImportIfNotExist(code, newPkg)

	assert.Contains(t, actual, newPkg)
	assert.Contains(t, actual, "fmt")

	t.Run("should not add if already exists", func(t *testing.T) {
		actual := code_editor.AddImportIfNotExist(code, "fmt")
		assert.Equal(t, actual, code)
	})
}
