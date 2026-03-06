package migration

import (
	"fmt"
	"time"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	name := c.Args().Get(0)
	if name == "" {
		return fmt.Errorf("migration name is required. Usage: catuaba g migration <name>")
	}

	timestamp := time.Now().Format("20060102150405")
	baseName := timestamp + "_" + generator.Snakeze(name)

	output.Info("Generating migration: %s", name)

	upPath := "database/migrations/" + baseName + ".up.sql"
	downPath := "database/migrations/" + baseName + ".down.sql"

	upContent := fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n", name, time.Now().Format(time.RFC3339))
	downContent := fmt.Sprintf("-- Rollback: %s\n-- Created at: %s\n\n", name, time.Now().Format(time.RFC3339))

	generator.GenerateFromContent(upContent, nil, upPath)
	generator.GenerateFromContent(downContent, nil, downPath)

	output.Success("Migration %s created!", name)
	output.Info("Next steps:")
	output.Info("  Edit the generated .up.sql and .down.sql files")
	output.Info("  Run 'catuaba db migrate' to apply migrations")

	return nil
}
