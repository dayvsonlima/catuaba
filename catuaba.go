package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "new",
				Aliases: []string{"n"},
				Usage:   "create new catuaba application",
				Action: func(c *cli.Context) error {
					fmt.Println("creates catuaba here")
					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "start web server",
				Action: func(c *cli.Context) error {
					fmt.Println("Start server here")
					return nil
				},
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
						Action: func(c *cli.Context) error {
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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
