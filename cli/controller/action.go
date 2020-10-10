package controller

import (
	"github.com/dayvsonlima/catuaba/cli/generator"
	"github.com/urfave/cli/v2"
)

type ControllerBuilder struct {
	Name       string
	MethodName string
}

func Action(c *cli.Context) error {
	var (
		name   = c.Args().Get(0)
		params = getAttributes(c)
	)

	generator.Mkdir(`app/controllers/` + generator.Snakeze(name))

	for _, methodName := range params {
		data := ControllerBuilder{
			Name:       c.Args().Get(0),
			MethodName: methodName,
		}

		controllerPath := "app/controllers/" + generator.Snakeze(name) + "/" + methodName + ".go"
		generator.GenerateFile("controller.go.tmpl", data, controllerPath)
	}

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
