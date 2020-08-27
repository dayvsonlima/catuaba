package scaffold

import (
	"strings"

	"github.com/dayvsonlima/catuaba/cli/code_editor"
	"github.com/dayvsonlima/catuaba/cli/generator"
	"github.com/dayvsonlima/catuaba/cli/model"
	"github.com/urfave/cli/v2"
)

type ControllerBuilder struct {
	Name       string
	MethodName string
	Params     []string
}

func Action(c *cli.Context) error {

	data := model.ModelBuilder{
		Name:   generator.Camelize(c.Args().Get(0)),
		Params: model.GetModelAttributes(c),
	}

	model.BuildModel(data)

	BuildControllers(data)
	BuildRoutes(data)
	return nil
}

func BuildControllers(data model.ModelBuilder) {
	name := data.Name
	controllerName := generator.Snakeze(generator.Pluralize(name))
	methods := []string{"create", "delete", "index", "show", "update"}

	generator.Mkdir(`app/controllers/` + controllerName)

	for _, methodName := range methods {
		data := ControllerBuilder{
			Name:       name,
			MethodName: methodName,
			Params:     data.Params,
		}

		controllerPath := "/app/controllers/" + controllerName + "/" + methodName + ".go"
		generator.GenerateFile("scaffold/controller/"+methodName+".go.tmpl", data, controllerPath)
	}

	generator.GenerateFile("scaffold/controller/shared.go.tmpl", data, "/app/controllers/"+controllerName+"/shared.go")
}

func BuildRoutes(data model.ModelBuilder) {

	routes := generator.Render("scaffold/routes.go.tmpl", data)

	code_editor.EditFile("/config/routes.go", func(content string) string {

		newPkg := "application/app/controllers/" + generator.Snakeze(generator.Pluralize(data.Name))

		routesString := strings.ReplaceAll(content, "\n}", routes+"\n}")
		routesString = code_editor.AddImport(routesString, newPkg)

		return routesString
	})
}
