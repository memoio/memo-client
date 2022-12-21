package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/memo-client/lib"
	"github.com/memoio/memo-client/lib/address"
	"github.com/memoio/memo-client/lib/repo"
	"github.com/memoio/memo-client/wallet"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/sha3"
)

var WalletCmd = &cli.Command{
	Name: "wallet",
	Subcommands: []*cli.Command{
		WalletListCmd,
		// WalletApproveCmd,
		// WalletDeposit,
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

var WalletApproveCmd = &cli.Command{
	Name:  "approve",
	Usage: "pay fee",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "taddr",
			Usage: "approve amount to addr",
		},
	},
	Action: func(cctx *cli.Context) error {
		taddr := cctx.String("taddr")

		buf, err := os.ReadFile("address")
		if err != nil {
			return err
		}

		bucket := string(buf)

		repoDir := cctx.String("repo")

		rep, err := repo.NewFSRepo(repoDir)
		if err != nil {
			return nil
		}

		defer func() {
			_ = rep.Close()
		}()

		pw := cctx.String("passwd")

		w := wallet.New(pw, rep.KeyStore())

		maddr := ethcommon.HexToAddress(bucket)
		if err != nil {
			return nil
		}

		srcaddr, err := address.NewAddress(maddr.Bytes())
		if err != nil {
			return err
		}

		// get sk
		sks, err := w.WalletExport(cctx.Context, srcaddr, "")
		if err != nil {
			return err
		}

		sk := hex.EncodeToString(sks.SecretKey)

		client, err := lib.New()
		if err != nil {
			return err
		}

		tokenaddr, err := client.GetTokenAddress(cctx.Context)
		if err != nil {
			return err
		}
		log.Println(tokenaddr)

		tokenaddress := ethcommon.HexToAddress(tokenaddr)
		toaddress := ethcommon.HexToAddress(taddr)

		value := big.NewInt(150503225806451)
		tshash, err := lib.Approve(cctx.Context, sk, tokenaddress, toaddress, value)
		if err != nil {
			return err
		}

		tsbyte, err := tshash.MarshalBinary()
		if err != nil {
			return err
		}
		ts := hex.EncodeToString(tsbyte)
		log.Println(ts)

		err = client.Approve(cctx.Context, ts, bucket)
		if err != nil {
			return err
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

var WalletDeposit = &cli.Command{
	Usage: "deposit",
	Name:  "deposit",
	Action: func(ctx *cli.Context) error {
		// buf, err := os.ReadFile("address")
		// if err != nil {
		// 	return err
		// }

		maddr := ethcommon.HexToAddress("0xecF059784977F181ECfc38C827ee818330bC76aD")

		client, err := ethclient.Dial("https://chain.metamemo.one:8501")
		if err != nil {
			return err
		}
		defer client.Close()
		storeaddr := ethcommon.HexToAddress("0x31e7829Ea2054fDF4BCB921eDD3a98a825242267")
		toaddr := ethcommon.HexToAddress("0xCcf7b7F747100f3393a75DDf6864589f76F4eA25")
		nonce, err := client.PendingNonceAt(ctx.Context, storeaddr)
		if err != nil {
			return err
		}
		log.Println("nonce: ", nonce)

		chainID, err := client.NetworkID(ctx.Context)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("chainID: ", chainID)

		depositFnSignature := []byte("storeDeposit(address,uint256,string)")
		hash := sha3.NewLegacyKeccak256()
		hash.Write(depositFnSignature)
		methodID := hash.Sum(nil)[:4]

		amount := big.NewInt(1000000000000000000)
		paddedAddress := ethcommon.LeftPadBytes(maddr.Bytes(), 32)
		paddedAmount := ethcommon.LeftPadBytes(amount.Bytes(), 32)
		paddedHashLen := ethcommon.LeftPadBytes(big.NewInt(int64(len(chainID.Bytes()))).Bytes(), 32)
		paddedHashOffset := ethcommon.LeftPadBytes(big.NewInt(32*3).Bytes(), 32)
		paddedChainID := ethcommon.RightPadBytes(chainID.Bytes(), 32)

		var data []byte
		data = append(data, methodID...)
		data = append(data, paddedAddress...)
		data = append(data, paddedAmount...)
		data = append(data, paddedHashOffset...)
		data = append(data, paddedHashLen...)
		data = append(data, paddedChainID...)

		gasLimit := uint64(300000)
		gasPrice := big.NewInt(1000)
		tx := types.NewTransaction(nonce, toaddr, big.NewInt(0), gasLimit, gasPrice, data)

		privateKey, err := crypto.HexToECDSA("8a87053d296a0f0b4600173773c8081b12917cef7419b2675943b0aa99429b62")
		if err != nil {
			log.Println("get privateKey error: ", err)
			return err
		}
		log.Println(privateKey)

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Println("signedTx error: ", err)
			return err
		}

		err = client.SendTransaction(ctx.Context, signedTx)
		if err != nil {
			log.Println(err)
			return err
		}

		log.Println("waiting tx complete...")
		time.Sleep(30 * time.Second)

		receipt, err := client.TransactionReceipt(ctx.Context, signedTx.Hash())
		if err != nil {
			return err
		}
		log.Println(receipt.Status)
		return nil
	},
}
