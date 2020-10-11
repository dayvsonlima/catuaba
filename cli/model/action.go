package model

import (
	"github.com/dayvsonlima/catuaba/code_editor"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/urfave/cli/v2"
)

type ModelBuilder struct {
	Name   string
	Params []string
}

func Action(c *cli.Context) error {

	data := ModelBuilder{
		Name:   generator.Camelize(c.Args().Get(0)),
		Params: GetModelAttributes(c),
	}

	BuildModel(data)

	return nil
}

func BuildModel(data ModelBuilder) {
	modelPath := "app/models/" + generator.Snakeze(data.Name) + ".go"

	generator.GenerateFile("model.go.tmpl", data, modelPath)

	code_editor.EditFile("database/connection.go", func(code string) string {
		return code_editor.InsertAttribute(code, "AutoMigrate", "&models."+data.Name+"{}")
	})
}

func GetModelAttributes(c *cli.Context) []string {
	l := c.Args().Len()
	var params []string

	for i := 1; i < l; i++ {
		params = append(params, c.Args().Get(i))
	}

	return params
}
