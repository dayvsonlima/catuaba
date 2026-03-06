package generator_test

import (
	"testing"

	"github.com/dayvsonlima/catuaba/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	t.Run("renders a template file with data", func(t *testing.T) {
		type data struct {
			Name       string
			MethodName string
		}
		result, err := generator.Render("controller.go.tmpl", data{
			Name:       "Post",
			MethodName: "index",
		})
		require.NoError(t, err)
		assert.Contains(t, result, "package post")
		assert.Contains(t, result, "func Index")
	})

	t.Run("renders model template with attributes", func(t *testing.T) {
		type data struct {
			Name   string
			Params []string
		}
		result, err := generator.Render("model.go.tmpl", data{
			Name:   "Post",
			Params: []string{"title:string", "body:string"},
		})
		require.NoError(t, err)
		assert.Contains(t, result, "type Post struct")
		assert.Contains(t, result, "Title string")
		assert.Contains(t, result, "Body string")
	})
}

func TestRenderFromContent(t *testing.T) {
	t.Run("renders inline template content", func(t *testing.T) {
		tmpl := "Hello {{.Name}}"
		type data struct{ Name string }
		result, err := generator.RenderFromContent(tmpl, data{Name: "World"})
		require.NoError(t, err)
		assert.Equal(t, "Hello World", result)
	})

	t.Run("renders with template functions", func(t *testing.T) {
		tmpl := "{{.Name | toSnake}}"
		type data struct{ Name string }
		result, err := generator.RenderFromContent(tmpl, data{Name: "MyModel"})
		require.NoError(t, err)
		assert.Equal(t, "my_model", result)
	})

	t.Run("renders with plural function", func(t *testing.T) {
		tmpl := "{{.Name | toPlural}}"
		type data struct{ Name string }
		result, err := generator.RenderFromContent(tmpl, data{Name: "post"})
		require.NoError(t, err)
		assert.Equal(t, "posts", result)
	})
}
