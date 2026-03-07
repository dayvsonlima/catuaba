package new

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func newTestApp() *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "new",
				Action: Action,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "db",
						Value: "postgres",
					},
				},
			},
		},
	}
}

func TestAction(t *testing.T) {
	t.Run("returns error when name is empty", func(t *testing.T) {
		err := newTestApp().Run([]string{"catuaba", "new"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "app name is required")
	})

	t.Run("returns error for invalid name", func(t *testing.T) {
		err := newTestApp().Run([]string{"catuaba", "new", "123invalid"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid app name")
	})

	t.Run("returns error for unsupported db driver", func(t *testing.T) {
		err := newTestApp().Run([]string{"catuaba", "new", "--db", "oracle", "myapp"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported database driver")
	})

	t.Run("validates name format with regex", func(t *testing.T) {
		assert.True(t, validName.MatchString("myapp"))
		assert.True(t, validName.MatchString("my-app"))
		assert.True(t, validName.MatchString("my_app"))
		assert.True(t, validName.MatchString("MyApp123"))
		assert.False(t, validName.MatchString("123app"))
		assert.False(t, validName.MatchString("my app"))
		assert.False(t, validName.MatchString(""))
	})

	t.Run("creates project structure with postgres", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(origDir)

		set := flag.NewFlagSet("test", flag.ContinueOnError)
		set.String("db", "postgres", "")
		set.Parse([]string{"testapp"})
		ctx := cli.NewContext(&cli.App{}, set, nil)

		err := Action(ctx)
		require.NoError(t, err)

		// Check directories
		assert.DirExists(t, filepath.Join(tmpDir, "testapp"))
		assert.DirExists(t, filepath.Join(tmpDir, "testapp", "config"))
		assert.DirExists(t, filepath.Join(tmpDir, "testapp", "database"))
		assert.DirExists(t, filepath.Join(tmpDir, "testapp", "database", "migrations"))
		assert.DirExists(t, filepath.Join(tmpDir, "testapp", "app", "controllers"))
		assert.DirExists(t, filepath.Join(tmpDir, "testapp", "app", "models"))
		assert.DirExists(t, filepath.Join(tmpDir, "testapp", "middleware"))

		// Check files
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "go.mod"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "application.go"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", ".env"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", ".env.example"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "config", "config.go"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "config", "routes.go"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "database", "connection.go"))

		// Check middleware files
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "middleware", "health.go"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "middleware", "logger.go"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "middleware", "cors.go"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "middleware", "recovery.go"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "middleware", "request_id.go"))

		// Check DevOps files
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "Dockerfile"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "docker-compose.yml"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "Makefile"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", ".gitignore"))
		assert.DirExists(t, filepath.Join(tmpDir, "testapp", ".github", "workflows"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", ".github", "workflows", "ci.yml"))

		// Check component files
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "nav.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "flash.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "pagination.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "form_field.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "badge.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "button.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "card.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "page_header.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "empty_state.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "back_link.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "detail_row.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "delete_form.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "table.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "spinner.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "form_actions.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "components", "error_page.templ"))
		assert.FileExists(t, filepath.Join(tmpDir, "testapp", "app", "views", "pages", "error.templ"))

		// Verify .env contains postgres config
		envContent, _ := os.ReadFile(filepath.Join(tmpDir, "testapp", ".env"))
		assert.Contains(t, string(envContent), "DB_DRIVER=postgres")
		assert.Contains(t, string(envContent), "DB_PORT=5432")

		// Verify application.go has graceful shutdown
		appContent, _ := os.ReadFile(filepath.Join(tmpDir, "testapp", "application.go"))
		assert.Contains(t, string(appContent), "signal.Notify")
		assert.Contains(t, string(appContent), "srv.Shutdown")

		// Verify docker-compose has postgres
		dcContent, _ := os.ReadFile(filepath.Join(tmpDir, "testapp", "docker-compose.yml"))
		assert.Contains(t, string(dcContent), "postgres:16-alpine")

		// Verify package.json has tailwind deps
		pkgContent, _ := os.ReadFile(filepath.Join(tmpDir, "testapp", "package.json"))
		assert.Contains(t, string(pkgContent), "tailwindcss")
		assert.Contains(t, string(pkgContent), "@tailwindcss/cli")
	})

	t.Run("creates project with sqlite", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(origDir)

		set := flag.NewFlagSet("test", flag.ContinueOnError)
		set.String("db", "sqlite", "")
		set.Parse([]string{"myapp"})
		ctx := cli.NewContext(&cli.App{}, set, nil)

		err := Action(ctx)
		require.NoError(t, err)

		envContent, _ := os.ReadFile(filepath.Join(tmpDir, "myapp", ".env"))
		assert.Contains(t, string(envContent), "DB_DRIVER=sqlite")
		assert.Contains(t, string(envContent), "DB_NAME=database/development.db")
	})
}
