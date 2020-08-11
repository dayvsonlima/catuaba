package model

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	// name := c.Args().Get(0)

	fmt.Println(c.Args().Get(0))
	fmt.Println(c.Args().Get(1))
	fmt.Println(c.Args().Get(2))
	fmt.Println(c.Args().Len())
	return nil
}
