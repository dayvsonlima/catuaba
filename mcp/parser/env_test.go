package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEnv(t *testing.T) {
	dir := t.TempDir()

	envContent := `# Application
APP_NAME=demo-app
APP_PORT=8080

# Database
DB_DRIVER=postgres
DB_HOST=localhost

# JWT
JWT_SECRET=mysecret
JWT_EXPIRY=24h
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte(envContent), 0644))

	vars := ParseEnv(dir)
	assert.Equal(t, []string{
		"APP_NAME", "APP_PORT",
		"DB_DRIVER", "DB_HOST",
		"JWT_SECRET", "JWT_EXPIRY",
	}, vars)
}

func TestParseEnv_FallbackToExample(t *testing.T) {
	dir := t.TempDir()

	envContent := `APP_NAME=test
APP_PORT=3000
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env.example"), []byte(envContent), 0644))

	vars := ParseEnv(dir)
	assert.Equal(t, []string{"APP_NAME", "APP_PORT"}, vars)
}

func TestParseEnv_NoFile(t *testing.T) {
	dir := t.TempDir()
	vars := ParseEnv(dir)
	assert.Nil(t, vars)
}
