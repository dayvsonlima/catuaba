package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func setupProject(t *testing.T) func() {
	t.Helper()
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	os.WriteFile("go.mod", []byte("module testapp\n"), 0644)
	os.MkdirAll("config", 0755)
	os.WriteFile("config/routes.go", []byte("package config\n"), 0644)
	return func() { os.Chdir(origDir) }
}

func TestAction(t *testing.T) {
	t.Run("returns error when name is empty", func(t *testing.T) {
		defer setupProject(t)()

		app := &cli.App{
			Commands: []*cli.Command{
				{Name: "service", Action: Action},
			},
		}
		err := app.Run([]string{"catuaba", "service"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "service name is required")
	})

	t.Run("generates service file", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(origDir)

		os.WriteFile("go.mod", []byte("module testapp\n"), 0644)
		os.MkdirAll("config", 0755)
		os.WriteFile("config/routes.go", []byte("package config\n"), 0644)

		app := &cli.App{
			Commands: []*cli.Command{
				{Name: "service", Action: Action},
			},
		}
		err := app.Run([]string{"catuaba", "service", "user"})
		require.NoError(t, err)

		assert.FileExists(t, filepath.Join(tmpDir, "app", "services", "user_service.go"))

		content, _ := os.ReadFile(filepath.Join(tmpDir, "app", "services", "user_service.go"))
		assert.Contains(t, string(content), "UserService")
		assert.Contains(t, string(content), "NewUserService")
		assert.Contains(t, string(content), "*gorm.DB")
	})
}
