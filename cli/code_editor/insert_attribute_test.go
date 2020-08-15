package code_editor_test

import (
	"testing"

	"github.com/dayvsonlima/catuaba/cli/code_editor"
	"github.com/stretchr/testify/assert"
)

var (
	code = `
		// Migrations .
		func Migrations() {
			Connection.AutoMigrate(&models.Post{})
		}
	`
)

func TestInsertAttribute(t *testing.T) {

	var (
		expectedOutput = "Connection.AutoMigrate(&models.Post{}, &models.User{})"
		methodName     = "AutoMigrate"
		newAttribute   = "&models.User{}"
	)

	out := code_editor.InsertAttribute(code, methodName, newAttribute)

	assert.Contains(t, out, expectedOutput)
}
