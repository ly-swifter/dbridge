package main

import (
	"github.com/urfave/cli/v2"
)

var sampleCmd = &cli.Command{
	Name:  "sample",
	Usage: "Initialize a dbridge node repo",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		return nil
	},
}
