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
