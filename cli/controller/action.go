package controller

import (
	"fmt"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/urfave/cli/v2"
)

type ControllerBuilder struct {
	Name       string
	MethodName string
}

func Action(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	name := c.Args().Get(0)
	if name == "" {
		return fmt.Errorf("controller name is required. Usage: catuaba g controller <name> [methods...]")
	}

	params := getAttributes(c)
	if len(params) == 0 {
		return fmt.Errorf("at least one method name is required. Usage: catuaba g controller <name> index show create")
	}

	output.Info("Generating controller: %s", name)

	if err := generator.Mkdir(`app/controllers/` + generator.Snakeze(name)); err != nil {
		return err
	}

	for _, methodName := range params {
		data := ControllerBuilder{
			Name:       c.Args().Get(0),
			MethodName: methodName,
		}

		controllerPath := "app/controllers/" + generator.Snakeze(name) + "/" + methodName + ".go"
		generator.GenerateFile("controller.go.tmpl", data, controllerPath)
	}

	output.Success("Controller %s generated successfully!", name)

	// Next steps
	snakeName := generator.Snakeze(name)
	output.Info("Next steps:")
	output.Info("  Register routes in config/routes.go:")
	output.Info("    routes.GET(\"/%s\", %s.%s)", snakeName, snakeName, generator.Camelize(params[0]))

	return nil
}

func getAttributes(c *cli.Context) []string {
	l := c.Args().Len()
	var params []string

	for i := 1; i < l; i++ {
		params = append(params, c.Args().Get(i))
	}

	return params
}
