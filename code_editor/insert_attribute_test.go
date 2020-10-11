package code_editor_test

import (
	"testing"

	"github.com/dayvsonlima/catuaba/code_editor"
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

		output := code_editor.InsertAttribute(code, methodName, newAttribute)
		assert.Contains(t, output, expectedOutput)
	})

	t.Run("when the method has no params", func(t *testing.T) {
		var (
			code = `
				// Migrations .
				func Migrations() {
					Connection.AutoMigrate( )
				}
			`

			expectedOutput = "Connection.AutoMigrate(&models.User{})"
		)

		output := code_editor.InsertAttribute(code, methodName, newAttribute)
		assert.Contains(t, output, expectedOutput)
	})

	t.Run("when the method has break lines", func(t *testing.T) {
		var (
			code = `
				// Migrations .
				func Migrations() {
					Connection.AutoMigrate(
						&models.Parms{},
						&models.Post{},
					)
				}
			`

			expectedOutput = "&models.User{}"
		)

		output := code_editor.InsertAttribute(code, methodName, newAttribute)
		assert.Contains(t, output, expectedOutput)
	})
}
