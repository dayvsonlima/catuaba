package seed

import (
	"fmt"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/urfave/cli/v2"
)

type SeedBuilder struct {
	Name string
}

func Action(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	name := c.Args().Get(0)
	if name == "" {
		return fmt.Errorf("seed name is required. Usage: catuaba g seed <name>")
	}

	data := SeedBuilder{
		Name: name,
	}

	filePath := "database/seeds/" + generator.Snakeze(name) + "_seed.go"

	output.Info("Generating seed: %s", name)
	if err := generator.Mkdir("database/seeds"); err != nil {
		return err
	}
	if err := generator.GenerateFile("seed.go.tmpl", data, filePath); err != nil {
		return err
	}
	output.Success("Seed %s generated!", name)
	output.Info("Next steps:")
	output.Info("  Call %sSeed(database.DB) from your seed runner", generator.Camelize(name))

	return nil
}
