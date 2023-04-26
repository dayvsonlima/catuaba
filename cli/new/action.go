package new

import (
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/dayvsonlima/catuaba/templates"
	"github.com/urfave/cli/v2"
)

type AppBuilder struct {
	Name   string
	Params []string
}

func Action(c *cli.Context) error {

	data := AppBuilder{
		Name: c.Args().Get(0),
	}

	// Root
	generator.Mkdir(data.Name)
	generator.GenerateFromContent(templates.Gomod, data, data.Name+"/go.mod")
	generator.GenerateFromContent(templates.Gosum, data, data.Name+"/go.sum")
	generator.GenerateFromContent(templates.Application, data, data.Name+"/application.go")

	// Config
	generator.Mkdir(data.Name + "/config")
	generator.GenerateFile("application/config/routes.go.tmpl", data, data.Name+"/config/routes.go")

	// Database
	generator.Mkdir(data.Name + "/database")
	generator.GenerateFile("application/database/connection.go.tmpl", data, data.Name+"/database/connection.go")

	// Controllers structure
	generator.Mkdir(data.Name + "/app")
	generator.Mkdir(data.Name + "/app/controllers")
	generator.GenerateFile("application/app/controllers/.keep", data, data.Name+"/app/controllers/.keep")

	// Models structure
	generator.Mkdir(data.Name + "/app/models")
	generator.GenerateFile("application/app/models/.keep", data, data.Name+"/app/models/.keep")

	return nil
}
