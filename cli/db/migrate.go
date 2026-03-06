package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func MigrateAction(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	db, driver, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := ensureMigrationsTable(db, driver); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to read applied migrations: %w", err)
	}

	files, err := getMigrationFiles("up")
	if err != nil {
		return err
	}

	pending := 0
	for _, f := range files {
		version := extractVersion(f)
		if applied[version] {
			continue
		}

		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", f, err)
		}

		sqlStr := strings.TrimSpace(string(content))
		if sqlStr == "" || !containsSQL(sqlStr) {
			output.Warning("Skipping empty migration: %s", filepath.Base(f))
			continue
		}

		output.Info("Running: %s", filepath.Base(f))
		if _, err := db.Exec(sqlStr); err != nil {
			return fmt.Errorf("migration %s failed: %w", filepath.Base(f), err)
		}

		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			// Retry with ? placeholder for MySQL/SQLite
			if _, err2 := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err2 != nil {
				return fmt.Errorf("failed to record migration %s: %w", version, err)
			}
		}

		pending++
	}

	if pending == 0 {
		output.Info("Nothing to migrate. Database is up to date.")
	} else {
		output.Success("Applied %d migration(s).", pending)
	}

	return nil
}

func RollbackAction(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	db, driver, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := ensureMigrationsTable(db, driver); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get the last applied migration
	var version string
	row := db.QueryRow("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1")
	if err := row.Scan(&version); err != nil {
		if err == sql.ErrNoRows {
			output.Info("Nothing to rollback.")
			return nil
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	// Find the corresponding .down.sql file
	downFile := ""
	files, _ := getMigrationFiles("down")
	for _, f := range files {
		if extractVersion(f) == version {
			downFile = f
			break
		}
	}

	if downFile == "" {
		return fmt.Errorf("rollback file not found for migration %s", version)
	}

	content, err := os.ReadFile(downFile)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", downFile, err)
	}

	sqlStr := strings.TrimSpace(string(content))
	if sqlStr != "" && containsSQL(sqlStr) {
		output.Info("Rolling back: %s", filepath.Base(downFile))
		if _, err := db.Exec(sqlStr); err != nil {
			return fmt.Errorf("rollback %s failed: %w", filepath.Base(downFile), err)
		}
	}

	// Remove from schema_migrations
	db.Exec("DELETE FROM schema_migrations WHERE version = $1", version)
	db.Exec("DELETE FROM schema_migrations WHERE version = ?", version)

	output.Success("Rolled back: %s", version)
	return nil
}

func StatusAction(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	db, driver, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := ensureMigrationsTable(db, driver); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to read applied migrations: %w", err)
	}

	files, err := getMigrationFiles("up")
	if err != nil {
		return err
	}

	if len(files) == 0 {
		output.Info("No migrations found in database/migrations/")
		return nil
	}

	for _, f := range files {
		version := extractVersion(f)
		status := "pending"
		if applied[version] {
			status = "applied"
		}
		name := filepath.Base(f)
		if status == "applied" {
			output.Success("  [applied]  %s", name)
		} else {
			output.Warning("  [pending]  %s", name)
		}
	}

	return nil
}

func connectDB() (*sql.DB, string, error) {
	godotenv.Load()

	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "postgres"
	}

	// Map user-facing driver names to Go sql driver names
	sqlDriver := driver
	switch driver {
	case "sqlite":
		sqlDriver = "sqlite3"
	}

	var dsn string
	switch driver {
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_USER", "postgres"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_NAME"),
			getEnv("DB_SSLMODE", "disable"),
		)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
			getEnv("DB_USER", "root"),
			os.Getenv("DB_PASS"),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "3306"),
			os.Getenv("DB_NAME"),
		)
	case "sqlite":
		dsn = getEnv("DB_NAME", "database/development.db")
	default:
		return nil, "", fmt.Errorf("unsupported database driver: %s", driver)
	}

	db, err := sql.Open(sqlDriver, dsn)
	if err != nil {
		return nil, "", fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, "", fmt.Errorf("failed to ping database: %w", err)
	}

	return db, driver, nil
}

func ensureMigrationsTable(db *sql.DB, driver string) error {
	query := `CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		applied[v] = true
	}
	return applied, nil
}

func getMigrationFiles(direction string) ([]string, error) {
	pattern := "database/migrations/*." + direction + ".sql"
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func extractVersion(filename string) string {
	base := filepath.Base(filename)
	// Format: 20240101120000_name.up.sql → 20240101120000
	parts := strings.SplitN(base, "_", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return base
}

func containsSQL(s string) bool {
	// Check if the content has actual SQL beyond just comments
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "--") {
			return true
		}
	}
	return false
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
