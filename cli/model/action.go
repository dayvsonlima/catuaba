package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/urfave/cli/v2"
)

type ModelBuilder struct {
	Name   string
	Params []string
}

func Action(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	name := c.Args().Get(0)
	if name == "" {
		return fmt.Errorf("model name is required. Usage: catuaba g model <name> [attributes...]")
	}

	attrs := GetModelAttributes(c)
	for _, attr := range attrs {
		if !strings.Contains(attr, ":") {
			return fmt.Errorf("invalid attribute format %q: expected name:type (e.g. title:string)", attr)
		}
	}

	data := ModelBuilder{
		Name:   generator.Camelize(name),
		Params: attrs,
	}

	return BuildModel(data)
}

func BuildModel(data ModelBuilder) error {
	modelPath := "app/models/" + generator.Snakeze(data.Name) + ".go"

	output.Info("Generating model: %s", data.Name)

	if err := generator.GenerateFile("model.go.tmpl", data, modelPath); err != nil {
		return err
	}

	// Generate SQL migration for the table
	if len(data.Params) > 0 {
		if err := generateMigration(data); err != nil {
			return err
		}
	}

	output.Success("Model %s generated successfully!", data.Name)
	return nil
}

func generateMigration(data ModelBuilder) error {
	timestamp := time.Now().Format("20060102150405")
	tableName := strings.ToLower(generator.Pluralize(generator.Snakeze(data.Name)))
	baseName := timestamp + "_create_" + tableName

	upPath := "database/migrations/" + baseName + ".up.sql"
	downPath := "database/migrations/" + baseName + ".down.sql"

	if err := generator.GenerateFile("migration_create_table.up.sql.tmpl", data, upPath); err != nil {
		return err
	}
	if err := generator.GenerateFile("migration_create_table.down.sql.tmpl", data, downPath); err != nil {
		return err
	}

	return nil
}

func GetModelAttributes(c *cli.Context) []string {
	l := c.Args().Len()
	var params []string

	for i := 1; i < l; i++ {
		params = append(params, c.Args().Get(i))
	}

	return params
}
