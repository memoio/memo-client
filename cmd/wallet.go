package cmd

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/memoio/memo-client/lib/repo"
	"github.com/memoio/memo-client/wallet"
	"github.com/urfave/cli/v2"
)

var WalletCmd = &cli.Command{
	Name: "wallet",
	Subcommands: []*cli.Command{
		WalletListCmd,
		// WalletSendCmd,
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
				toAddress := ethcommon.BytesToAddress(as.Bytes())
				fmt.Println(toAddress)
			}
		}

		return nil
	},
}

// var WalletSendCmd = &cli.Command{
// 	Name:      "send",
// 	Usage:     "send memo to another address",
// 	UsageText: "send [destaddress][value]",
// 	Action: func(ctx *cli.Context) error {
// 		args := ctx.Args()
// 		destaddr := args.Get(0)
// 		log.Println(destaddr)
// 		if destaddr == "" {
// 			return xerrors.New("addr is nil ")
// 		}

// 		value := args.Get(1)
// 		if value == "" {
// 			return xerrors.New("value is nil")
// 		}
// 		// ivalue, err := strconv.ParseInt(value, 10, 64)
// 		// if err != nil {
// 		// 	return err
// 		// }

// 		repoDir := ctx.String("repo")

// 		rep, err := repo.NewFSRepo(repoDir)
// 		if err != nil {
// 			return nil
// 		}

// 		defer func() {
// 			_ = rep.Close()
// 		}()

// 		pw := ctx.String("passwd")

// 		buf, err := os.ReadFile("address")
// 		if err != nil {
// 			return err
// 		}

// 		maddr := ethcommon.HexToAddress(string(buf))
// 		if err != nil {
// 			return nil
// 		}

// 		srcaddr, err := address.NewAddress(maddr.Bytes())
// 		if err != nil {
// 			return err
// 		}

// 		sk, err := wallet.GetSk(ctx.Context, repoDir, pw, srcaddr)
// 		if err != nil {
// 			return err
// 		}

// 		hash, err := lib.Signmsg(ctx.Context, sk, value, destaddr, "")
// 		if err != nil {
// 			return err
// 		}
// 		log.Println(hash)
// 		return nil
// 	},
// }
