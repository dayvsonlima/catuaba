package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertLineAfter(t *testing.T) {
	code := `package main

import "fmt"

func main() {
	fmt.Println("hello")
}`

	result := InsertLineAfter(code, `fmt.Println("hello")`, `	fmt.Println("world")`)
	assert.Contains(t, result, "hello\")\n\tfmt.Println(\"world\")")
}

func TestInsertLineAfter_NotFound(t *testing.T) {
	code := "package main\n"
	result := InsertLineAfter(code, "nonexistent", "new line")
	assert.Equal(t, code, result)
}

func TestInsertBeforeClosingBrace(t *testing.T) {
	code := `func SetupRoutes(r *gin.Engine) {
	r.GET("/health", health.Check)
}`

	result := InsertBeforeClosingBrace(code, "\tr.POST(\"/auth/login\", auth.Login)")
	assert.Contains(t, result, "auth.Login)")
	assert.True(t, len(result) > len(code))
	// The closing brace should still be present
	assert.Contains(t, result, "}")
}

func TestAppendToFile(t *testing.T) {
	code := "JWT_SECRET=old\n"
	result := AppendToFile(code, "NEW_VAR=value")
	assert.Equal(t, "JWT_SECRET=old\nNEW_VAR=value\n", result)
}

func TestAppendToFile_NoTrailingNewline(t *testing.T) {
	code := "JWT_SECRET=old"
	result := AppendToFile(code, "NEW_VAR=value")
	assert.Equal(t, "JWT_SECRET=old\nNEW_VAR=value\n", result)
}
