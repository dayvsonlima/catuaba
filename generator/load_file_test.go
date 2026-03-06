package generator_test

import (
	"testing"

	"github.com/dayvsonlima/catuaba/generator"
	"github.com/stretchr/testify/assert"
)

func TestLoadFile(t *testing.T) {
	t.Run("loads an existing template file", func(t *testing.T) {
		content := generator.LoadFile("controller.go.tmpl")
		assert.NotEmpty(t, content)
		assert.Contains(t, content, "package")
	})

	t.Run("loads a nested template file", func(t *testing.T) {
		content := generator.LoadFile("application/application.go.tmpl")
		assert.NotEmpty(t, content)
		assert.Contains(t, content, "gin.New()")
	})

	t.Run("returns empty string for non-existent file", func(t *testing.T) {
		content := generator.LoadFile("non_existent.tmpl")
		assert.Empty(t, content)
	})
}
