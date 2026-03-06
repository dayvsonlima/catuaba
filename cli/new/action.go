package new

import (
	"fmt"
	"os"
	"regexp"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/urfave/cli/v2"
)

type AppBuilder struct {
	Name     string
	DBDriver string
	DBHost   string
	DBPort   string
	DBUser   string
	DBName   string
	Auth     bool
}

var validName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

var dbDefaults = map[string]AppBuilder{
	"postgres": {DBHost: "localhost", DBPort: "5432", DBUser: "postgres"},
	"mysql":    {DBHost: "localhost", DBPort: "3306", DBUser: "root"},
	"sqlite":   {DBHost: "", DBPort: "", DBUser: ""},
}

func Action(c *cli.Context) error {
	name := c.Args().Get(0)
	if name == "" {
		return fmt.Errorf("app name is required. Usage: catuaba new <name>")
	}
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid app name %q: must start with a letter and contain only letters, numbers, hyphens and underscores", name)
	}

	dbDriver := c.String("db")
	defaults, ok := dbDefaults[dbDriver]
	if !ok {
		return fmt.Errorf("unsupported database driver %q: use postgres, mysql, or sqlite", dbDriver)
	}

	dbName := name + "_development"
	if dbDriver == "sqlite" {
		dbName = "database/development.db"
	}

	data := AppBuilder{
		Name:     name,
		DBDriver: dbDriver,
		DBHost:   defaults.DBHost,
		DBPort:   defaults.DBPort,
		DBUser:   defaults.DBUser,
		DBName:   dbName,
		Auth:     c.Bool("auth"),
	}

	// Warn if directory already exists
	if info, err := os.Stat(name); err == nil && info.IsDir() {
		output.Warning("Directory %q already exists. Files may be overwritten.", name)
	}

	authLabel := ""
	if data.Auth {
		authLabel = ", auth: enabled"
	}
	output.Info("Creating new Catuaba application: %s (db: %s%s)", name, dbDriver, authLabel)

	dirs := []string{
		data.Name,
		data.Name + "/config",
		data.Name + "/database",
		data.Name + "/database/migrations",
		data.Name + "/app",
		data.Name + "/app/controllers",
		data.Name + "/app/models",
		data.Name + "/app/views/layouts",
		data.Name + "/app/views/components",
		data.Name + "/app/views/pages",
		data.Name + "/middleware",
		data.Name + "/static/css",
		data.Name + "/static/js",
		data.Name + "/.github",
		data.Name + "/.github/workflows",
	}

	if data.Auth {
		dirs = append(dirs,
			data.Name+"/app/controllers/auth",
			data.Name+"/app/views/auth",
		)
	}

	for _, dir := range dirs {
		if err := generator.Mkdir(dir); err != nil {
			return err
		}
	}

	files := []struct {
		tmpl string
		dest string
	}{
		{"application/go.mod.tmpl", data.Name + "/go.mod"},
		{"application/application.go.tmpl", data.Name + "/application.go"},
		{"application/dot-env.tmpl", data.Name + "/.env"},
		{"application/dot-env-example.tmpl", data.Name + "/.env.example"},
		{"application/config/config.go.tmpl", data.Name + "/config/config.go"},
		{"application/config/routes.go.tmpl", data.Name + "/config/routes.go"},
		{"application/database/connection.go.tmpl", data.Name + "/database/connection.go"},
		{"application/middleware/health.go.tmpl", data.Name + "/middleware/health.go"},
		{"application/middleware/logger.go.tmpl", data.Name + "/middleware/logger.go"},
		{"application/middleware/cors.go.tmpl", data.Name + "/middleware/cors.go"},
		{"application/middleware/recovery.go.tmpl", data.Name + "/middleware/recovery.go"},
		{"application/middleware/request_id.go.tmpl", data.Name + "/middleware/request_id.go"},
		{"application/middleware/flash.go.tmpl", data.Name + "/middleware/flash.go"},
		{"application/middleware/csrf.go.tmpl", data.Name + "/middleware/csrf.go"},
		{"application/middleware/session.go.tmpl", data.Name + "/middleware/session.go"},
		{"application/middleware/rate_limit.go.tmpl", data.Name + "/middleware/rate_limit.go"},
		{"application/middleware/secure_headers.go.tmpl", data.Name + "/middleware/secure_headers.go"},
		{"application/Dockerfile.tmpl", data.Name + "/Dockerfile"},
		{"application/docker-compose.yml.tmpl", data.Name + "/docker-compose.yml"},
		{"application/Makefile.tmpl", data.Name + "/Makefile"},
		{"application/dot-gitignore.tmpl", data.Name + "/.gitignore"},
		{"application/github-actions.yml.tmpl", data.Name + "/.github/workflows/ci.yml"},
		{"application/air.toml.tmpl", data.Name + "/.air.toml"},
		{"application/README.md.tmpl", data.Name + "/README.md"},
		{"application/CLAUDE.md.tmpl", data.Name + "/CLAUDE.md"},
		{"application/tailwind.config.js.tmpl", data.Name + "/tailwind.config.js"},
		{"application/static/css/input.css.tmpl", data.Name + "/static/css/input.css"},
		{"application/controllers/home.go.tmpl", data.Name + "/app/controllers/home.go"},
		{"application/views/layouts/base.templ.tmpl", data.Name + "/app/views/layouts/base.templ"},
		{"application/views/components/nav.templ.tmpl", data.Name + "/app/views/components/nav.templ"},
		{"application/views/components/flash.templ.tmpl", data.Name + "/app/views/components/flash.templ"},
		{"application/views/components/pagination.templ.tmpl", data.Name + "/app/views/components/pagination.templ"},
		{"application/views/components/form_field.templ.tmpl", data.Name + "/app/views/components/form_field.templ"},
		{"application/views/pages/home.templ.tmpl", data.Name + "/app/views/pages/home.templ"},
		{"application/views/pages/not_found.templ.tmpl", data.Name + "/app/views/pages/not_found.templ"},
	}

	if data.Auth {
		files = append(files,
			struct{ tmpl, dest string }{"application/models/user.go.tmpl", data.Name + "/app/models/user.go"},
			struct{ tmpl, dest string }{"application/controllers/auth/login.go.tmpl", data.Name + "/app/controllers/auth/login.go"},
			struct{ tmpl, dest string }{"application/controllers/auth/register.go.tmpl", data.Name + "/app/controllers/auth/register.go"},
			struct{ tmpl, dest string }{"application/controllers/auth/logout.go.tmpl", data.Name + "/app/controllers/auth/logout.go"},
			struct{ tmpl, dest string }{"application/views/auth/login.templ.tmpl", data.Name + "/app/views/auth/login.templ"},
			struct{ tmpl, dest string }{"application/views/auth/register.templ.tmpl", data.Name + "/app/views/auth/register.templ"},
			struct{ tmpl, dest string }{"application/middleware/require_auth.go.tmpl", data.Name + "/middleware/require_auth.go"},
			struct{ tmpl, dest string }{"migration_create_users.up.sql.tmpl", data.Name + "/database/migrations/00000000000001_create_users.up.sql"},
			struct{ tmpl, dest string }{"migration_create_users.down.sql.tmpl", data.Name + "/database/migrations/00000000000001_create_users.down.sql"},
		)
	}

	for _, f := range files {
		if err := generator.GenerateFile(f.tmpl, data, f.dest); err != nil {
			return err
		}
	}

	// Keep files for empty directories
	keepFiles := []string{
		data.Name + "/database/migrations/.keep",
		data.Name + "/app/controllers/.keep",
		data.Name + "/app/models/.keep",
		data.Name + "/static/js/.keep",
	}
	for _, kf := range keepFiles {
		if err := generator.GenerateFromContent("", data, kf); err != nil {
			return err
		}
	}

	output.Success("Application %s created successfully!", name)
	output.Info("Next steps:")
	output.Info("  cd %s", name)
	output.Info("  go mod tidy")
	output.Info("  go install github.com/a-h/templ/cmd/templ@latest")
	output.Info("  make dev")

	return nil
}
