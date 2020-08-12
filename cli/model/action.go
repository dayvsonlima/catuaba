package model

import (
	"fmt"

	"github.com/dayvsonlima/catuaba/cli/templates"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	data := struct {
		Name string
	}{
		Name: c.Args().Get(0),
	}

	out := templates.Render("cli/model/template.go.tmpl", data)

	fmt.Println(out)
	return nil
}
