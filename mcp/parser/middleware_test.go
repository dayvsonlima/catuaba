package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMiddleware(t *testing.T) {
	dir := t.TempDir()

	authMW := `package middleware

import "net/http"

func RequireAuth(next http.Handler) http.Handler {
	return nil
}
`
	corsMW := `package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
	return nil
}
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "auth.go"), []byte(authMW), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "cors.go"), []byte(corsMW), 0644))

	middlewares, err := ParseMiddleware(dir)
	require.NoError(t, err)
	assert.Len(t, middlewares, 2)

	names := make(map[string]string)
	for _, m := range middlewares {
		names[m.Name] = m.File
	}
	assert.Equal(t, "auth.go", names["RequireAuth"])
	assert.Equal(t, "cors.go", names["CORS"])
}

func TestParseMiddleware_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	middlewares, err := ParseMiddleware(dir)
	require.NoError(t, err)
	assert.Empty(t, middlewares)
}
