package cmd

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/memo-client/lib/repo"
	"github.com/memoio/memo-client/wallet"
	"github.com/urfave/cli/v2"
)

var WalletCmd = &cli.Command{
	Name: "wallet",
	Subcommands: []*cli.Command{
		WalletListCmd,
	},
}

var WalletListCmd = &cli.Command{
	Name:  "list",
	Usage: "list wallet address",
	Action: func(ctx *cli.Context) error {
		repoDir := ctx.String("repo")

		rep, err := repo.NewFSRepo(repoDir)
		if err != nil {
			return nil
		}

		defer func() {
			_ = rep.Close()
		}()

		pw := ctx.String("passwd")

		w := wallet.New(pw, rep.KeyStore())
		addrs, err := w.WalletList(ctx.Context)
		if err != nil {
			return err
		}
		
		for _, as := range addrs {
			if as.Len() == 20 {
				toAddress := common.BytesToAddress(as.Bytes())
				fmt.Println(toAddress)
			}
		}

		return nil
	},
}
