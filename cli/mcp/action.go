package mcp

import (
	"os"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/generator"
	catuabamcp "github.com/dayvsonlima/catuaba/mcp"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	if err := generator.IsInsideCatuabaProject(); err != nil {
		return err
	}

	// Redirect all CLI output to stderr so stdout stays clean for JSON-RPC
	output.SetWriter(os.Stderr)

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	return catuabamcp.Serve(projectDir)
}
