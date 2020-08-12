package model

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dayvsonlima/catuaba/cli/templates"
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

	// out := templates.Render("cli/model/template.go.tmpl", data)
	out := templates.RenderFromContent(Template, data)

	currentPath, _ := os.Getwd()
	modelPath := currentPath + "/app/models/" + templates.Snakeze(data.Name) + ".go"
	err := ioutil.WriteFile(modelPath, []byte(out), 0644)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}

	fmt.Println("+" + modelPath)
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
