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
	assert.Equal(t, "ConversationID", generator.Camelize("conversation_id"))
	assert.Equal(t, "UserID", generator.Camelize("user_id"))
	assert.Equal(t, "APIURL", generator.Camelize("api_url"))
	assert.Equal(t, "HTTPURL", generator.Camelize("http_url"))
	assert.Equal(t, "ID", generator.Camelize("id"))
	assert.Equal(t, "DB", generator.Camelize("db"))
}

func TestSnakeze(t *testing.T) {
	assert.Equal(t, "my_model", generator.Snakeze("MyModel"))
	assert.Equal(t, "user_profile", generator.Snakeze("UserProfile"))
	assert.Equal(t, "html_parser", generator.Snakeze("HTMLParser"))
	assert.Equal(t, "conversation_id", generator.Snakeze("ConversationID"))
	assert.Equal(t, "user_id", generator.Snakeze("UserID"))
	assert.Equal(t, "httpurl", generator.Snakeze("HTTPURL"))   // consecutive acronyms merge (edge case)
	assert.Equal(t, "http_server", generator.Snakeze("HTTPServer"))
}

func TestLowerPlural(t *testing.T) {
	assert.Equal(t, "posts", generator.LowerPlural("Post"))
	assert.Equal(t, "categories", generator.LowerPlural("Category"))
}

func TestGetAttributeName(t *testing.T) {
	assert.Equal(t, "Title", generator.GetAttributeName("title:string"))
	assert.Equal(t, "UserName", generator.GetAttributeName("user_name:string"))
	assert.Equal(t, "ConversationID", generator.GetAttributeName("conversation_id:integer"))
	assert.Equal(t, "UserID", generator.GetAttributeName("user_id:integer"))
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
	assert.Equal(t, "post", generator.CamelizeVar("post"))
	assert.Equal(t, "mySuperVarName", generator.CamelizeVar("my_super_var_name"))
	assert.Equal(t, "conversationID", generator.CamelizeVar("conversation_id"))
	assert.Equal(t, "userID", generator.CamelizeVar("user_id"))
	assert.Equal(t, "id", generator.CamelizeVar("id"))
	assert.Equal(t, "apiURL", generator.CamelizeVar("api_url"))
}
