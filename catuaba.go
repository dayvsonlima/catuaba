package main

import (
	"os"
	"sort"

	"github.com/dayvsonlima/catuaba/cli/api"
	"github.com/dayvsonlima/catuaba/cli/controller"
	clidb "github.com/dayvsonlima/catuaba/cli/db"
	"github.com/dayvsonlima/catuaba/cli/install"
	climcp "github.com/dayvsonlima/catuaba/cli/mcp"
	"github.com/dayvsonlima/catuaba/cli/middleware"
	"github.com/dayvsonlima/catuaba/cli/migration"
	"github.com/dayvsonlima/catuaba/cli/model"
	"github.com/dayvsonlima/catuaba/cli/new"
	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/cli/plugincmd"
	"github.com/dayvsonlima/catuaba/cli/routes"
	"github.com/dayvsonlima/catuaba/cli/scaffold"
	"github.com/dayvsonlima/catuaba/cli/seed"
	"github.com/dayvsonlima/catuaba/cli/server"
	"github.com/dayvsonlima/catuaba/cli/service"
	"github.com/dayvsonlima/catuaba/cli/upgrade"
	"github.com/dayvsonlima/catuaba/cli/view"

	"github.com/urfave/cli/v2"
)

func main() {
	version := "0.1.9"
	app := &cli.App{
		Name:     "catuaba",
		Usage:    "Build full-stack Go web apps in minutes, not days",
		Version:  version,
		Metadata: map[string]interface{}{"version": version},
		Commands: []*cli.Command{
			{
				Name:   "mcp",
				Usage:  "start MCP server (stdio) for AI integration",
				Action: climcp.Action,
			},
			{
				Name:      "new",
				Aliases:   []string{"n"},
				Usage:     "create new catuaba application",
				ArgsUsage: "<app-name>",
				Action:    new.Action,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "db",
						Value: "postgres",
						Usage: "database driver (postgres, mysql, sqlite)",
					},
					&cli.BoolFlag{
						Name:  "auth",
						Usage: "include authentication (login, register, logout)",
					},
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "start web server",
				Action:  server.Action,
			},
			{
				Name:    "routes",
				Aliases: []string{"r"},
				Usage:   "list all application routes",
				Action:  routes.Action,
			},
			{
				Name:  "db",
				Usage: "database commands (migrate, rollback, status)",
				Subcommands: []*cli.Command{
					{
						Name:    "migrate",
						Aliases: []string{"m"},
						Usage:   "run pending migrations",
						Action:  clidb.MigrateAction,
					},
					{
						Name:    "rollback",
						Aliases: []string{"r"},
						Usage:   "rollback the last migration",
						Action:  clidb.RollbackAction,
					},
					{
						Name:    "status",
						Aliases: []string{"s"},
						Usage:   "show migration status",
						Action:  clidb.StatusAction,
					},
				},
				Action: func(c *cli.Context) error {
					output.Info("Usage: catuaba db <command>")
					return cli.ShowSubcommandHelp(c)
				},
			},
			{
				Name:      "install",
				Aliases:   []string{"i"},
				Usage:     "install a plugin",
				ArgsUsage: "<name|url|path>",
				Action:    install.Action,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:  "var",
						Usage: "set plugin variable (key=value)",
					},
				},
			},
			{
				Name:   "upgrade",
				Usage:  "upgrade catuaba to the latest version",
				Action: upgrade.Action,
			},
			{
				Name:  "plugin",
				Usage: "manage plugins",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "list available plugins from registry",
						Action: plugincmd.ListAction,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "refresh",
								Usage: "force refresh registry cache",
							},
						},
					},
					{
						Name:      "info",
						Usage:     "show plugin details",
						ArgsUsage: "<name>",
						Action:    plugincmd.InfoAction,
					},
				},
				Action: func(c *cli.Context) error {
					output.Info("Usage: catuaba plugin <command>")
					return cli.ShowSubcommandHelp(c)
				},
			},
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "generate code (controllers, models, scaffolds, etc.)",
				Subcommands: []*cli.Command{
					{
						Name:      "controller",
						Aliases:   []string{"c"},
						Usage:     "generate controller",
						ArgsUsage: "<name> <method> [method...]",
						Action:    controller.Action,
					},
					{
						Name:      "model",
						Aliases:   []string{"m"},
						Usage:     "generate model",
						ArgsUsage: "<name> [field:type...]",
						Action:    model.Action,
					},
					{
						Name:      "scaffold",
						Aliases:   []string{"s"},
						Usage:     "generate full-stack scaffold (model + HTML handlers + views + routes)",
						ArgsUsage: "<name> [field:type...]",
						Action:    scaffold.Action,
					},
					{
						Name:      "api",
						Aliases:   []string{"a"},
						Usage:     "generate JSON API (model + controller + routes)",
						ArgsUsage: "<name> [field:type...]",
						Action:    api.Action,
					},
					{
						Name:      "migration",
						Aliases:   []string{"mi"},
						Usage:     "generate migration files",
						ArgsUsage: "<name>",
						Action:    migration.Action,
					},
					{
						Name:      "middleware",
						Aliases:   []string{"mw"},
						Usage:     "generate middleware",
						ArgsUsage: "<name>",
						Action:    middleware.Action,
					},
					{
						Name:      "service",
						Aliases:   []string{"sv"},
						Usage:     "generate service",
						ArgsUsage: "<name>",
						Action:    service.Action,
					},
					{
						Name:      "seed",
						Aliases:   []string{"sd"},
						Usage:     "generate database seed",
						ArgsUsage: "<name>",
						Action:    seed.Action,
					},
					{
						Name:      "view",
						Aliases:   []string{"v"},
						Usage:     "generate a view page",
						ArgsUsage: "<name>",
						Action:    view.Action,
					},
				},
				Action: func(c *cli.Context) error {
					output.Info("Usage: catuaba generate <command> [options]")
					return cli.ShowSubcommandHelp(c)
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.EnableBashCompletion = true
	if err := app.Run(os.Args); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
}
