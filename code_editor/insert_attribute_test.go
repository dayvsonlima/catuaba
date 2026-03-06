package code_editor_test

import (
	"testing"

	"github.com/dayvsonlima/catuaba/code_editor"
	"github.com/stretchr/testify/assert"
)

func TestInsertAttribute(t *testing.T) {

	t.Run("when the method has some param (AutoMigrate)", func(t *testing.T) {
		var (
			methodName     = "AutoMigrate"
			newAttribute   = "&models.User{}"
			code           = `
			func Migrations() {
				Connection.AutoMigrate(&models.Post{})
			}
			`
			expectedOutput = "Connection.AutoMigrate(&models.Post{}, &models.User{})"
		)

		output := code_editor.InsertAttribute(code, methodName, newAttribute)
		assert.Contains(t, output, expectedOutput)
	})

	t.Run("when the method has no params (AutoMigrate)", func(t *testing.T) {
		var (
			methodName     = "AutoMigrate"
			newAttribute   = "&models.User{}"
			code           = `
				func Migrations() {
					Connection.AutoMigrate( )
				}
			`
			expectedOutput = "Connection.AutoMigrate(&models.User{})"
		)

		output := code_editor.InsertAttribute(code, methodName, newAttribute)
		assert.Contains(t, output, expectedOutput)
	})

	t.Run("when the method has break lines (AutoMigrate)", func(t *testing.T) {
		var (
			methodName     = "AutoMigrate"
			newAttribute   = "&models.User{}"
			code           = `
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

	t.Run("works with Migrate method (GORM v2)", func(t *testing.T) {
		var (
			methodName     = "Migrate"
			newAttribute   = "&models.User{}"
			code           = `
			func init() {
				database.Migrate(&models.Post{})
			}
			`
			expectedOutput = "database.Migrate(&models.Post{}, &models.User{})"
		)

		output := code_editor.InsertAttribute(code, methodName, newAttribute)
		assert.Contains(t, output, expectedOutput)
	})

	t.Run("works with empty Migrate (GORM v2)", func(t *testing.T) {
		var (
			methodName     = "Migrate"
			newAttribute   = "&models.User{}"
			code           = `
			func init() {
				database.Migrate( )
			}
			`
			expectedOutput = "database.Migrate(&models.User{})"
		)

		output := code_editor.InsertAttribute(code, methodName, newAttribute)
		assert.Contains(t, output, expectedOutput)
	})
}
