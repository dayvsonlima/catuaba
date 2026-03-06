package install

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/generator"
	"github.com/dayvsonlima/catuaba/plugin"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	source := c.Args().Get(0)
	if source == "" {
		return fmt.Errorf("plugin name, URL, or path is required.\nUsage: catuaba install <name|url|path> [--var key=value]")
	}

	// Parse --var flags
	vars := parseVars(c.StringSlice("var"))

	output.Info("Resolving plugin: %s", source)

	// 1. Resolve source to local dir
	pluginDir, cleanup, err := plugin.ResolveSource(source)
	if err != nil {
		return err
	}
	defer cleanup()

	// 2. Load manifest
	m, err := plugin.LoadManifest(pluginDir + "/plugin.yaml")
	if err != nil {
		return err
	}

	// 3. Validate
	if err := plugin.Validate(m); err != nil {
		return err
	}

	// 4. Check if already installed
	tracker, err := plugin.LoadTracker()
	if err == nil {
		if _, installed := tracker.Installed[m.Name]; installed {
			output.Warning("Plugin %q is already installed. Reinstalling...", m.Name)
		}
	}

	// 5. Check file conflicts
	conflicts := plugin.CheckConflicts(m)
	if len(conflicts) > 0 {
		output.Warning("The following files already exist and will be overwritten:")
		for _, f := range conflicts {
			output.Warning("  %s", f)
		}
	}

	output.Info("Installing plugin: %s v%s", m.Name, m.Version)
	fmt.Println()

	// 6. Run installation
	if err := plugin.Install(pluginDir, m, vars); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// 7. Run go get for dependencies
	if len(m.Dependencies) > 0 {
		output.Info("Installing Go dependencies...")
		for _, dep := range m.Dependencies {
			cmd := exec.Command("go", "get", dep)
			cmd.Dir = "."
			if out, err := cmd.CombinedOutput(); err != nil {
				output.Warning("Failed to install %s: %s", dep, string(out))
			} else {
				output.Info("  go get %s", dep)
			}
		}
	}

	// 8. Track installation
	if err := plugin.TrackInstall(m); err != nil {
		output.Warning("Failed to save plugin tracking: %v", err)
	}

	// 9. Post-install messages
	fmt.Println()
	output.Success("Plugin %s installed successfully!", m.Name)
	if len(m.PostInstall) > 0 {
		fmt.Println()
		output.Info("Post-install notes:")
		for _, msg := range m.PostInstall {
			output.Info("  • %s", msg)
		}
	}

	return nil
}

func parseVars(flags []string) map[string]string {
	vars := make(map[string]string)
	for _, f := range flags {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars
}
