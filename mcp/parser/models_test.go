package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseModels(t *testing.T) {
	dir := t.TempDir()

	userModel := `package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string ` + "`json:\"email\" gorm:\"uniqueIndex;not null\"`" + `
	PasswordHash string ` + "`json:\"-\"`" + `
	Name         string ` + "`json:\"name\"`" + `
}
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "user.go"), []byte(userModel), 0644))

	productModel := `package models

type Product struct {
	Title string ` + "`json:\"title\"`" + `
	Price float64 ` + "`json:\"price\"`" + `
}
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "product.go"), []byte(productModel), 0644))

	// Should be ignored
	require.NoError(t, os.WriteFile(filepath.Join(dir, "user_test.go"), []byte("package models"), 0644))

	models, err := ParseModels(dir)
	require.NoError(t, err)
	assert.Len(t, models, 2)

	// Find User model
	var user, product = models[1], models[0]
	if models[0].Name == "User" {
		user, product = models[0], models[1]
	}

	assert.Equal(t, "User", user.Name)
	assert.Len(t, user.Fields, 4)
	assert.Equal(t, "gorm.Model", user.Fields[0].Name)
	assert.Equal(t, "embedded", user.Fields[0].Type)
	assert.Equal(t, "Email", user.Fields[1].Name)
	assert.Equal(t, "string", user.Fields[1].Type)
	assert.Equal(t, "email", user.Fields[1].JSON)
	assert.Equal(t, "uniqueIndex;not null", user.Fields[1].GORM)

	assert.Equal(t, "Product", product.Name)
	assert.Len(t, product.Fields, 2)
}

func TestParseModels_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	models, err := ParseModels(dir)
	require.NoError(t, err)
	assert.Empty(t, models)
}

func TestParseModels_NoDir(t *testing.T) {
	_, err := ParseModels("/nonexistent/dir")
	assert.Error(t, err)
}
