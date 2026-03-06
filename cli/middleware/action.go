package middleware

import (
	"fmt"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/urfave/cli/v2"
)

type MiddlewareBuilder struct {
	Name string
}

func Action(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	name := c.Args().Get(0)
	if name == "" {
		return fmt.Errorf("middleware name is required. Usage: catuaba g middleware <name>")
	}

	data := MiddlewareBuilder{
		Name: name,
	}

	filePath := "middleware/" + generator.Snakeze(name) + ".go"

	output.Info("Generating middleware: %s", name)
	if err := generator.GenerateFile("middleware.go.tmpl", data, filePath); err != nil {
		return err
	}
	output.Success("Middleware %s generated!", name)
	output.Info("Next steps:")
	output.Info("  Add to application.go: engine.Use(middleware.%s())", generator.Camelize(name))

	return nil
}
