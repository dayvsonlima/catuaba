package plugincmd

import (
	"fmt"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/plugin"
	"github.com/urfave/cli/v2"
)

func ListAction(c *cli.Context) error {
	refresh := c.Bool("refresh")

	reg, err := plugin.LoadRegistry(refresh)
	if err != nil {
		return fmt.Errorf("loading registry: %w", err)
	}

	if len(reg.Plugins) == 0 {
		output.Info("No plugins available in the registry.")
		return nil
	}

	output.Info("Available plugins:\n")
	for name, entry := range reg.Plugins {
		fmt.Printf("  %-15s %s (v%s)\n", name, entry.Description, entry.Version)
	}
	fmt.Println()
	output.Info("Install with: catuaba install <name>")

	return nil
}

func InfoAction(c *cli.Context) error {
	name := c.Args().Get(0)
	if name == "" {
		return fmt.Errorf("plugin name is required. Usage: catuaba plugin info <name>")
	}

	reg, err := plugin.LoadRegistry(false)
	if err != nil {
		return fmt.Errorf("loading registry: %w", err)
	}

	entry, ok := reg.Plugins[name]
	if !ok {
		return fmt.Errorf("plugin %q not found in registry", name)
	}

	fmt.Printf("Name:        %s\n", name)
	fmt.Printf("Version:     %s\n", entry.Version)
	fmt.Printf("Description: %s\n", entry.Description)
	fmt.Printf("Repository:  %s\n", entry.Repository)
	fmt.Println()
	output.Info("Install with: catuaba install %s", name)

	return nil
}
