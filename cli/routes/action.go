package routes

import (
	"fmt"
	"os"
	"strings"

	"github.com/dayvsonlima/catuaba/generator"
	"github.com/dayvsonlima/catuaba/mcp/parser"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var methodColors = map[string]*color.Color{
	"GET":     color.New(color.FgGreen),
	"POST":    color.New(color.FgYellow),
	"PUT":     color.New(color.FgBlue),
	"PATCH":   color.New(color.FgCyan),
	"DELETE":  color.New(color.FgRed),
	"HEAD":    color.New(color.FgWhite),
	"OPTIONS": color.New(color.FgWhite),
}

func Action(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	var allRoutes []parser.RouteInfo

	// Parse application.go for top-level routes (home, auth, static)
	if _, err := os.Stat("application.go"); err == nil {
		if routes, err := parser.ParseRoutes("application.go"); err == nil {
			allRoutes = append(allRoutes, routes...)
		}
	}

	// Parse config/routes.go for scaffold/api routes
	if _, err := os.Stat("config/routes.go"); err == nil {
		if routes, err := parser.ParseRoutes("config/routes.go"); err == nil {
			allRoutes = append(allRoutes, routes...)
		}
	}

	if len(allRoutes) == 0 {
		fmt.Println("  No routes found.")
		return nil
	}

	// Calculate column widths
	maxPath := 4
	for _, r := range allRoutes {
		if len(r.Path) > maxPath {
			maxPath = len(r.Path)
		}
	}

	// Print table
	header := color.New(color.Faint)
	header.Printf("  %-7s  %-*s  %s\n", "METHOD", maxPath, "PATH", "HANDLER")
	fmt.Printf("  %s  %s  %s\n", strings.Repeat("─", 7), strings.Repeat("─", maxPath), strings.Repeat("─", 30))

	for _, r := range allRoutes {
		mc, ok := methodColors[r.Method]
		if !ok {
			mc = color.New(color.FgWhite)
		}
		fmt.Printf("  %s  %-*s  %s\n", mc.Sprintf("%-7s", r.Method), maxPath, r.Path, r.Handler)
	}

	fmt.Printf("\n  %d routes total\n", len(allRoutes))
	return nil
}
