package server

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	output.Info("Starting development server...")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cmd := exec.Command("air")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	go func() {
		<-sigs
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	if err := cmd.Run(); err != nil {
		output.Error("Error running air: %v", err)
		output.Info("Make sure 'air' is installed: go install github.com/air-verse/air@latest")
		return err
	}

	return nil
}
