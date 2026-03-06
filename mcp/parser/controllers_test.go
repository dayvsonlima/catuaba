package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseControllers(t *testing.T) {
	dir := t.TempDir()

	// Create auth controller package
	authDir := filepath.Join(dir, "auth")
	require.NoError(t, os.MkdirAll(authDir, 0755))

	registerGo := `package auth

import "net/http"

func Register(w http.ResponseWriter, r *http.Request) {}
`
	loginGo := `package auth

import "net/http"

func Login(w http.ResponseWriter, r *http.Request) {}

// helper is unexported and should be skipped
func helper() {}
`
	require.NoError(t, os.WriteFile(filepath.Join(authDir, "register.go"), []byte(registerGo), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(authDir, "login.go"), []byte(loginGo), 0644))

	// Should be ignored
	require.NoError(t, os.WriteFile(filepath.Join(authDir, "register_test.go"), []byte("package auth"), 0644))

	controllers, err := ParseControllers(dir)
	require.NoError(t, err)
	assert.Len(t, controllers, 1)
	assert.Equal(t, "auth", controllers[0].Package)
	assert.Len(t, controllers[0].Functions, 2)

	names := make(map[string]string)
	for _, f := range controllers[0].Functions {
		names[f.Name] = f.File
	}
	assert.Equal(t, "login.go", names["Login"])
	assert.Equal(t, "register.go", names["Register"])
}

func TestParseControllers_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	controllers, err := ParseControllers(dir)
	require.NoError(t, err)
	assert.Empty(t, controllers)
}
