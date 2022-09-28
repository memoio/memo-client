package cmd

import (
	"context"
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/memo-client/lib/crypto/signature"
	"github.com/memoio/memo-client/lib/repo"
	"github.com/memoio/memo-client/lib/types"
	"github.com/memoio/memo-client/lib/utils"
	"github.com/memoio/memo-client/wallet"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"lukechampine.com/blake3"
)

var InitCmd = &cli.Command{
	Name:  "init",
	Usage: "init a memo client",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "repo",
			Usage: "key repo",
			Value: "./",
		},
	},
	Action: func(ctx *cli.Context) error {
		log.Println("Initializing memo client")
		repoDir := ctx.String("repo")

		exist, err := repo.Exists(repoDir)
		if err != nil {
			return err
		}

		if exist {
			return xerrors.Errorf("repo at '%s' is already initialized", repoDir)
		}

		log.Println("initializing repo at: ", repoDir)

		rep, err := repo.NewFSRepo(repoDir)
		if err != nil {
			return err
		}

		defer func() {
			_ = rep.Close()
		}()
		pw := ctx.String("passwd")
		sk := ctx.String("sk")
		err = create(ctx.Context, rep, pw, sk)

		if err != nil {
			log.Printf ("fail initializing node %s", err)
			return err
		}
		return nil
	},
}

func create(ctx context.Context, r repo.Repo, password, sk string) error {
	w := wallet.New(password, r.KeyStore())

	var sBytes []byte
	if sk == "" {
		log.Println("generating wallet address...")

		privkey, err := signature.GenerateKey(types.Secp256k1)
		if err != nil {
			return err
		}

		sbytes, err := privkey.Raw()
		if err != nil {
			return nil
		}

		sBytes = sbytes
	} else {
		sbytes, err := hex.DecodeString(sk)
		if err != nil {
			return err
		}

		sBytes = sbytes
	}

	wki := &types.KeyInfo{
		Type:      types.Secp256k1,
		SecretKey: sBytes,
	}

	addr, err := w.WalletImport(ctx, wki)
	if err != nil {
		return err
	}

	wa := common.BytesToAddress(utils.ToEthAddress(addr.Bytes()))

	if sk == "" {
		log.Println("generated wallet address: ", wa)
	} else {
		log.Println("import wallet address: ", wa)
	}

	log.Println("generating bls key...")

	blsSeed := make([]byte, len(sBytes)+1)
	copy(blsSeed[:len(sBytes)], sBytes)
	blsSeed[len(sBytes)] = byte(types.BLS)
	blsByte := blake3.Sum256(blsSeed)
	blsKey := &types.KeyInfo{
		SecretKey: blsByte[:],
		Type:      types.BLS,
	}

	blsAddr, err := w.WalletImport(ctx, blsKey)
	if err != nil {
		return err
	}

	log.Println("genenrated bls key: ", blsAddr.String())

	return nil

}
