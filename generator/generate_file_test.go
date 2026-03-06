package generator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dayvsonlima/catuaba/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFile(t *testing.T) {
	t.Run("generates a file from a template", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(origDir)

		type data struct {
			Name       string
			MethodName string
		}

		err := generator.GenerateFile("controller.go.tmpl", data{
			Name:       "Post",
			MethodName: "index",
		}, "index.go")

		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(tmpDir, "index.go"))
		require.NoError(t, err)
		assert.Contains(t, string(content), "package post")
		assert.Contains(t, string(content), "func Index")
	})
}

func TestGenerateFromContent(t *testing.T) {
	t.Run("generates a file from inline content", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(origDir)

		type data struct{ Name string }
		err := generator.GenerateFromContent("Hello {{.Name}}", data{Name: "World"}, "hello.txt")

		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(tmpDir, "hello.txt"))
		require.NoError(t, err)
		assert.Equal(t, "Hello World", string(content))
	})

	t.Run("generates empty file from empty content", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(origDir)

		type data struct{}
		err := generator.GenerateFromContent("", data{}, ".keep")

		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(tmpDir, ".keep"))
		require.NoError(t, err)
		assert.Empty(t, string(content))
	})
}
