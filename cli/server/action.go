package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cosmtrek/air/runner"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	r, err := runner.NewEngine("", false)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	go func() {
		<-sigs
		r.Stop()
	}()

	defer func() {
		if e := recover(); e != nil {
			log.Fatalf("PANIC: %+v", e)
		}
	}()

	r.Run()
	return nil
}
