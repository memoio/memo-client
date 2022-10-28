package cmd

import (
	"fmt"
	"os"

	"github.com/memoio/memo-client/lib"
	"github.com/urfave/cli/v2"
)

var QueryCmd = &cli.Command{
	Name:  "query",
	Usage: "query price ",
	Action: func(ctx *cli.Context) error {
		client, err := lib.New()
		if err != nil {
			return err
		}

		price, err := client.QueryPrice(ctx.Context)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(price)

		return nil
	},
}

var GetBalanceInfoCmd = &cli.Command{
	Name:  "balance",
	Usage: "get balance info",
	Action: func(ctx *cli.Context) error {
		client, err := lib.New()
		if err != nil {
			return err
		}

		buf, err := os.ReadFile("address")
		if err != nil {
			return err
		}

		address := string(buf)
		balance, err := client.GetBalanceInfo(ctx.Context, address)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(balance)

		return nil
	},
}
