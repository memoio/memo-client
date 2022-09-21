package main

import (
	"fmt"
	"os"

	"github.com/memoio/memo-client/cmd"
	"github.com/memoio/memo-client/lib"
	"github.com/urfave/cli/v2"
)

func main() {
	ak := os.Getenv("ACCESS_KEY")
	if ak == "" {
		ak = "memo"
	}

	sk := os.Getenv("SECRET_KEY")
	if sk == "" {
		sk = "memoriae"
	}

	endpoint := os.Getenv("ENDPOINT")
	if endpoint == "" {
		endpoint = "0.0.0.0:5080"
	}

	lib.AccessKey = ak
	lib.SecretKey = sk
	lib.EndPoint = endpoint

	local := make([]*cli.Command, 0)
	local = append(local, cmd.PutObjectCmd)
	local = append(local, cmd.QueryCmd)
	local = append(local, cmd.GetBalanceInfoCmd)

	app := &cli.App{
		Name:  "memo-client",
		Usage: "memo client",

		Commands: local,
	}

	app.Setup()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err)
		os.Exit(1)
	}
}
