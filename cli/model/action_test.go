package model

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestActionValidation(t *testing.T) {
	t.Run("returns error when name is empty", func(t *testing.T) {
		defer setupProject(t)()

		app := &cli.App{
			Commands: []*cli.Command{
				{Name: "model", Action: Action},
			},
		}
		err := app.Run([]string{"catuaba", "model"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "model name is required")
	})

	t.Run("returns error for invalid attribute format", func(t *testing.T) {
		defer setupProject(t)()

		app := &cli.App{
			Commands: []*cli.Command{
				{Name: "model", Action: Action},
			},
		}
		err := app.Run([]string{"catuaba", "model", "Post", "title"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid attribute format")
	})
}
