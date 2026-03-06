package scaffold

import (
	"fmt"
	"strings"

	"github.com/dayvsonlima/catuaba/cli/model"
	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/code_editor"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/urfave/cli/v2"
)

type HandlerBuilder struct {
	Name       string
	MethodName string
	Params     []string
}

func Action(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	name := c.Args().Get(0)
	if name == "" {
		return fmt.Errorf("scaffold name is required. Usage: catuaba g scaffold <name> [attributes...]")
	}

	attrs := model.GetModelAttributes(c)
	for _, attr := range attrs {
		if !strings.Contains(attr, ":") {
			return fmt.Errorf("invalid attribute format %q: expected name:type (e.g. title:string)", attr)
		}
	}

	data := model.ModelBuilder{
		Name:   generator.Camelize(name),
		Params: attrs,
	}

	if err := model.BuildModel(data); err != nil {
		return err
	}

	output.Info("Generating scaffold for: %s", data.Name)
	if err := BuildHandlers(data); err != nil {
		return err
	}
	if err := BuildViews(data); err != nil {
		return err
	}
	if err := BuildRoutes(data); err != nil {
		return err
	}
	output.Success("Scaffold %s generated successfully!", data.Name)

	routeBase := "/" + generator.Snakeze(generator.Pluralize(data.Name))
	output.Info("Routes created:")
	output.Route("GET", routeBase)
	output.Route("GET", routeBase+"/new")
	output.Route("POST", routeBase)
	output.Route("GET", routeBase+"/:id")
	output.Route("GET", routeBase+"/:id/edit")
	output.Route("POST", routeBase+"/:id")
	output.Route("POST", routeBase+"/:id/delete")
	output.Info("Next steps:")
	output.Info("  Run 'make dev' to start the server")
	output.Info("  Visit: http://localhost:8080%s", routeBase)

	return nil
}

func BuildHandlers(data model.ModelBuilder) error {
	name := data.Name
	handlerName := generator.Snakeze(generator.Pluralize(name))
	methods := []string{"index", "show", "new", "create", "edit", "update", "delete"}

	if err := generator.Mkdir("app/controllers/" + handlerName); err != nil {
		return err
	}

	for _, methodName := range methods {
		handlerData := HandlerBuilder{
			Name:       name,
			MethodName: methodName,
			Params:     data.Params,
		}

		handlerPath := "app/controllers/" + handlerName + "/" + methodName + ".go"
		if err := generator.GenerateFile("scaffold/handler/"+methodName+".go.tmpl", handlerData, handlerPath); err != nil {
			return err
		}
	}

	return nil
}

func BuildViews(data model.ModelBuilder) error {
	name := data.Name
	viewName := generator.Snakeze(generator.Pluralize(name))
	views := []string{"index", "show", "form"}

	if err := generator.Mkdir("app/views/" + viewName); err != nil {
		return err
	}

	for _, viewFile := range views {
		viewPath := "app/views/" + viewName + "/" + viewFile + ".templ"
		if err := generator.GenerateFile("scaffold/view/"+viewFile+".templ.tmpl", data, viewPath); err != nil {
			return err
		}
	}

	return nil
}

func BuildRoutes(data model.ModelBuilder) error {
	routes, err := generator.Render("scaffold/view_routes.go.tmpl", data)
	if err != nil {
		return err
	}

	return code_editor.EditFile("config/routes.go", func(content string) string {
		moduleName := generator.ModuleName()
		newPkg := moduleName + "/app/controllers/" + generator.Snakeze(generator.Pluralize(data.Name))

		marker := "// [catuaba:routes]"
		routesString := strings.Replace(content, marker, routes+marker, 1)
		routesString = code_editor.AddImport(routesString, newPkg)

		return routesString
	})
}
