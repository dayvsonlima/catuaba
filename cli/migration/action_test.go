package migration

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

func TestActionValidation(t *testing.T) {
	t.Run("returns error when name is empty", func(t *testing.T) {
		defer setupProject(t)()

		app := &cli.App{
			Commands: []*cli.Command{
				{Name: "migration", Action: Action},
			},
		}
		err := app.Run([]string{"catuaba", "migration"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "migration name is required")
	})

	t.Run("creates migration files", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(origDir)

		os.WriteFile("go.mod", []byte("module testapp\n"), 0644)
		os.MkdirAll("config", 0755)
		os.WriteFile("config/routes.go", []byte("package config\n"), 0644)
		os.MkdirAll("database/migrations", 0755)

		app := &cli.App{
			Commands: []*cli.Command{
				{Name: "migration", Action: Action},
			},
		}
		err := app.Run([]string{"catuaba", "migration", "create_users"})
		require.NoError(t, err)

		entries, _ := os.ReadDir(filepath.Join(tmpDir, "database", "migrations"))
		sqlFiles := 0
		hasUp := false
		hasDown := false
		for _, e := range entries {
			if e.Name() == ".keep" {
				continue
			}
			sqlFiles++
			if filepath.Ext(e.Name()) == ".sql" {
				if assert.Contains(t, e.Name(), "create_users") {
					if len(e.Name()) > 0 {
						if filepath.Ext(e.Name()[:len(e.Name())-4]) == ".up" {
							hasUp = true
						}
						if filepath.Ext(e.Name()[:len(e.Name())-4]) == ".down" {
							hasDown = true
						}
					}
				}
			}
		}
		assert.Equal(t, 2, sqlFiles)
		assert.True(t, hasUp, "should have .up.sql file")
		assert.True(t, hasDown, "should have .down.sql file")
	})
}
