package generator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dayvsonlima/catuaba/generator"
)

func TestPluralize(t *testing.T) {
	assert.Equal(t, "posts", generator.Pluralize("post"))
	assert.Equal(t, "categories", generator.Pluralize("category"))
	assert.Equal(t, "people", generator.Pluralize("person"))
}

func TestCamelize(t *testing.T) {
	assert.Equal(t, "Post", generator.Camelize("post"))
	assert.Equal(t, "MyModel", generator.Camelize("my_model"))
	assert.Equal(t, "UserProfile", generator.Camelize("user_profile"))
}

func TestSnakeze(t *testing.T) {
	assert.Equal(t, "my_model", generator.Snakeze("MyModel"))
	assert.Equal(t, "user_profile", generator.Snakeze("UserProfile"))
	assert.Equal(t, "html_parser", generator.Snakeze("HTMLParser"))
}

func TestLowerPlural(t *testing.T) {
	assert.Equal(t, "posts", generator.LowerPlural("Post"))
	assert.Equal(t, "categories", generator.LowerPlural("Category"))
}

func TestGetAttributeName(t *testing.T) {
	assert.Equal(t, "Title", generator.GetAttributeName("title:string"))
	assert.Equal(t, "UserName", generator.GetAttributeName("user_name:string"))
}

func TestGetAttributeType(t *testing.T) {
	assert.Equal(t, "string", generator.GetAttributeType("title:string"))
	assert.Equal(t, "int", generator.GetAttributeType("age:int"))
}

func TestGetAttributeJson(t *testing.T) {
	result := generator.GetAttributeJson("user_name:string")
	assert.Contains(t, result, `json:"user_name"`)
	assert.NotContains(t, result, `binding:"required"`)
}

func TestGetAttributeJsonBinding(t *testing.T) {
	result := generator.GetAttributeJsonBinding("user_name:string")
	assert.Contains(t, result, `json:"user_name"`)
	assert.Contains(t, result, `binding:"required"`)
}

func TestCamelizeVar(t *testing.T) {

	t.Run("When receive a single text", func(t *testing.T) {
		expected := "post"
		output := generator.CamelizeVar("Post")
		assert.Equal(t, expected, output)
	})

	t.Run("When receive a complex text", func(t *testing.T) {
		expected := "mySuperVarName"
		output := generator.CamelizeVar("MySuper_var_name")
		assert.Equal(t, expected, output)
	})
}
