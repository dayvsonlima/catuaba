package view

import (
	"fmt"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/urfave/cli/v2"
)

type ViewBuilder struct {
	Name string
}

func Action(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	name := c.Args().Get(0)
	if name == "" {
		return fmt.Errorf("view name is required. Usage: catuaba g view <name>")
	}

	data := ViewBuilder{
		Name: generator.Camelize(name),
	}

	viewPath := "app/views/pages/" + generator.Snakeze(name) + ".templ"
	if err := generator.GenerateFile("view.templ.tmpl", data, viewPath); err != nil {
		return err
	}

	output.Success("View %s generated!", data.Name)
	output.Info("  File: %s", viewPath)
	output.Info("  Add a route in application.go or config/routes.go to serve it")

	return nil
}
