package cmd

import "github.com/urfave/cli/v2"

var initCmd = &cli.Command{
	Name:  "init",
	Usage: "init a memo client",
	Action: func(ctx *cli.Context) error {
		
		return nil
	},
}
