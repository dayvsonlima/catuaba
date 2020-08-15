package code_editor_test

import (
	"testing"

	"github.com/dayvsonlima/catuaba/cli/code_editor"
	"github.com/stretchr/testify/assert"
)

func TestInsertAttribute(t *testing.T) {

	var (
		methodName   = "AutoMigrate"
		newAttribute = "&models.User{}"
	)

	t.Run("when the method has some param", func(t *testing.T) {
		var (
			code = `
			// Migrations .
			func Migrations() {
				Connection.AutoMigrate(&models.Post{})
			}
			`
			expectedOutput = "Connection.AutoMigrate(&models.Post{}, &models.User{})"
		)

		out := code_editor.InsertAttribute(code, methodName, newAttribute)
		assert.Contains(t, out, expectedOutput)
	})

	t.Run("when the method has no params", func(t *testing.T) {
		var (
			code = `
				// Migrations .
				func Migrations() {
					Connection.AutoMigrate()
				}
			`

			expectedOutput = "Connection.AutoMigrate(&models.Post{})"
		)

		out := code_editor.InsertAttribute(code, methodName, newAttribute)
		assert.Contains(t, out, expectedOutput)
	})
}
