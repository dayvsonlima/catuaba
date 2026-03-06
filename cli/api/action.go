package api

import (
	"fmt"
	"strings"

	"github.com/dayvsonlima/catuaba/cli/model"
	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/code_editor"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/urfave/cli/v2"
)

type ControllerBuilder struct {
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
		return fmt.Errorf("api name is required. Usage: catuaba g api <name> [attributes...]")
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

	output.Info("Generating API controllers for: %s", data.Name)
	if err := BuildControllers(data); err != nil {
		return err
	}
	if err := BuildRoutes(data); err != nil {
		return err
	}
	output.Success("API %s generated successfully!", data.Name)

	routeBase := "/" + generator.Snakeze(generator.Pluralize(data.Name))
	output.Info("Routes created:")
	output.Route("GET", "/api"+routeBase)
	output.Route("POST", "/api"+routeBase)
	output.Route("GET", "/api"+routeBase+"/:id")
	output.Route("PUT", "/api"+routeBase+"/:id")
	output.Route("PATCH", "/api"+routeBase+"/:id")
	output.Route("DELETE", "/api"+routeBase+"/:id")
	output.Info("Next steps:")
	output.Info("  Run 'make run' to start the server")
	output.Info("  Try: curl http://localhost:8080/api%s", routeBase)

	return nil
}

func BuildControllers(data model.ModelBuilder) error {
	name := data.Name
	controllerName := generator.Snakeze(generator.Pluralize(name))
	methods := []string{"create", "delete", "index", "show", "update"}

	if err := generator.Mkdir("app/controllers/api/" + controllerName); err != nil {
		return err
	}

	for _, methodName := range methods {
		controllerData := ControllerBuilder{
			Name:       name,
			MethodName: methodName,
			Params:     data.Params,
		}

		controllerPath := "app/controllers/api/" + controllerName + "/" + methodName + ".go"
		generator.GenerateFile("scaffold/controller/"+methodName+".go.tmpl", controllerData, controllerPath)

		testPath := "app/controllers/api/" + controllerName + "/" + methodName + "_test.go"
		generator.GenerateFile("scaffold/controller/"+methodName+"_test.go.tmpl", controllerData, testPath)
	}

	return nil
}

func BuildRoutes(data model.ModelBuilder) error {
	routes, err := generator.Render("scaffold/api_routes.go.tmpl", data)
	if err != nil {
		return err
	}

	return code_editor.EditFile("config/routes.go", func(content string) string {
		moduleName := generator.ModuleName()
		resourcePkg := generator.Snakeze(generator.Pluralize(data.Name))
		newPkg := moduleName + "/app/controllers/api/" + resourcePkg
		alias := "api_" + resourcePkg

		routesString := strings.ReplaceAll(content, "\n}", routes+"\n}")
		routesString = code_editor.AddAliasedImport(routesString, alias, newPkg)

		return routesString
	})
}
