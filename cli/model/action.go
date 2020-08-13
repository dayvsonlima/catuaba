package model

import (
	"github.com/dayvsonlima/catuaba/cli/generator"
	"github.com/urfave/cli/v2"
)

type ModelBuilder struct {
	Name   string
	Params []string
}

func Action(c *cli.Context) error {

	data := ModelBuilder{
		Name:   c.Args().Get(0),
		Params: getAttributes(c),
	}

	modelPath := "/app/models/" + generator.Snakeze(data.Name) + ".go"

	generator.GenerateFile("model.go.tmpl", data, modelPath)

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
