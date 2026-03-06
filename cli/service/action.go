package service

import (
	"fmt"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/urfave/cli/v2"
)

type ServiceBuilder struct {
	Name string
}

func Action(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	name := c.Args().Get(0)
	if name == "" {
		return fmt.Errorf("service name is required. Usage: catuaba g service <name>")
	}

	data := ServiceBuilder{
		Name: name,
	}

	filePath := "app/services/" + generator.Snakeze(name) + "_service.go"

	output.Info("Generating service: %s", name)
	if err := generator.Mkdir("app/services"); err != nil {
		return err
	}
	if err := generator.GenerateFile("service.go.tmpl", data, filePath); err != nil {
		return err
	}
	output.Success("Service %s generated!", name)
	output.Info("Next steps:")
	output.Info("  Inject: svc := services.New%sService(database.DB)", generator.Camelize(name))

	return nil
}
