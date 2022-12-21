package cmd

import (
	"fmt"
	"math/big"
	"os"

	"github.com/memoio/memo-client/lib"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

var QueryCmd = &cli.Command{
	Name:  "query",
	Usage: "query price ",
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name:  "time",
			Usage: "time to storage(day)(min=100, max=1000)",
			Value: 100,
		},
		&cli.StringFlag{
			Name:  "path",
			Usage: "path of file",
		},
	},
	Action: func(cctx *cli.Context) error {
		client, err := lib.New()
		if err != nil {
			return err
		}
		buf, err := os.ReadFile("address")
		if err != nil {
			return err
		}

		bucket := string(buf)

		path := cctx.String("path")
		if path == "" {
			return xerrors.New("path is nil")
		}

		fileinfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		date := cctx.Int64("time")
		if date > 1000 || date < 100 {
			return xerrors.Errorf("time too long or too short")
		}
		time := big.NewInt(date)

		size := big.NewInt(fileinfo.Size())

		price, err := client.QueryPrice(cctx.Context, bucket, size.String(), time.String())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s automemo\n", price)

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
		fmt.Printf("%s automemo\n", balance)

		return nil
	},
}
