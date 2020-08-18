package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/dayvsonlima/catuaba/cli/controller"
	"github.com/dayvsonlima/catuaba/cli/new"
	"github.com/dayvsonlima/catuaba/cli/server"

	"github.com/dayvsonlima/catuaba/cli/model"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "new",
				Aliases: []string{"n"},
				Usage:   "create new catuaba application",
				Action:  new.Action,
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "start web server",
				Action:  server.Action,
			},
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "generate code structs",
				Subcommands: []*cli.Command{
					{
						Name:    "controller",
						Aliases: []string{"c"},
						Usage:   "generate controller",
						Action:  controller.Action,
					},
					{
						Name:    "model",
						Aliases: []string{"m"},
						Usage:   "generate model",
						Action:  model.Action,
					},
					{
						Name:    "scaffold",
						Aliases: []string{"s"},
						Usage:   "generate scaffold",
						Action: func(c *cli.Context) error {
							fmt.Println("generete your scaffold here")
							return nil
						},
					},
				},

				Action: func(c *cli.Context) error {
					fmt.Println("my awesome command")
					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.EnableBashCompletion = true
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
