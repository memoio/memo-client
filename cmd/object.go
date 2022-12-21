package cmd

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/memo-client/lib"
	"github.com/memoio/memo-client/lib/address"
	"github.com/memoio/memo-client/lib/repo"
	"github.com/memoio/memo-client/wallet"
	miniogo "github.com/minio/minio-go/v7"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

var PutObjectCmd = &cli.Command{
	Name:  "put",
	Usage: "put object",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "path",
			Usage:    "path of file",
			Required: true,
		},
		&cli.Int64Flag{
			Name:  "time",
			Usage: "time to storage(day)(min=100, max=1000)",
			Value: 100,
		},
	},
	Action: func(cctx *cli.Context) error {
		// get parameters
		buf, err := os.ReadFile("address")
		if err != nil {
			return err
		}

		bucket := string(buf)
		path := cctx.String("path")
		if path == "" {
			return xerrors.New("path is nil")
		}

		date := cctx.Int64("time")
		if date > 1000 || date < 100 {
			return xerrors.Errorf("time too long or too short")
		}

		dated := big.NewInt(date)

		fileinfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		ssize := big.NewInt(fileinfo.Size())

		client, err := lib.New()
		if err != nil {
			return err
		}

		price, err := client.QueryPrice(cctx.Context, bucket, ssize.String(), dated.String())
		if err != nil {
			fmt.Println(err)
		}

		amount := new(big.Int)
		amount.SetString(price, 10)

		balance, err := client.GetBalanceInfo(cctx.Context, bucket)
		if err != nil {
			return err
		}
		balancei := new(big.Int)
		balancei.SetString(balance, 10)

		if balancei.Cmp(amount) < 0 {
			return xerrors.Errorf("balance not enough, amount: %d, balance: %d", amount, balancei)
		}

		log.Printf("upload info: size is %dB, time is %dday, cost is %d automemo\n", fileinfo.Size(), date, amount)

		// ask whether to upload
		upload := false

		for i := 0; i < 3; i++ {
			res, err := ask4confirm()
			if err == nil {
				upload = res
				break
			}
		}

		if !upload {
			log.Println("cancel upload")
			return nil
		}

		repoDir := cctx.String("repo")

		rep, err := repo.NewFSRepo(repoDir)
		if err != nil {
			return nil
		}

		defer func() {
			_ = rep.Close()
		}()

		pw := cctx.String("passwd")

		maddr := ethcommon.HexToAddress(bucket)
		if err != nil {
			return nil
		}

		srcaddr, err := address.NewAddress(maddr.Bytes())
		if err != nil {
			return err
		}
		log.Println("srcaddr ", srcaddr)

		w := wallet.New(pw, rep.KeyStore())

		// get sk
		sks, err := w.WalletExport(cctx.Context, srcaddr, "")
		if err != nil {
			return err
		}

		log.Printf("%x\n", string(sks.SecretKey))
		sk := hex.EncodeToString(sks.SecretKey)

		privateKey, err := crypto.HexToECDSA(sk)
		if err != nil {
			log.Fatal(err)
		}

		sdata := []byte(dated.String())
		hash := crypto.Keccak256Hash(sdata)
		log.Println(hash.Hex())

		signature, err := crypto.Sign(hash.Bytes(), privateKey)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("sign: ", hexutil.Encode(signature))

		log.Println(fileinfo.Name())
		object := fileinfo.Name()
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		r := bytes.NewBuffer(data)

		metadata := make(map[string]string)

		metadata["sign"] = hexutil.Encode(signature)
		metadata["date"] = dated.String()

		log.Println("metadata: ", metadata)

		opt := miniogo.PutObjectOptions{
			UserMetadata:     metadata,
			DisableMultipart: true,
		}

		info, err := client.PutObject(cctx.Context, bucket, object, r, fileinfo.Size(), opt)
		if err != nil {
			return err
		}

		fmt.Println("cid Info:", info.ETag)
		return nil
	},
}

var GetObjectCmd = &cli.Command{
	Name:  "get",
	Usage: "get object",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "object",
			Aliases: []string{"on"},
			Usage:   "objectName",
		},
		&cli.StringFlag{
			Name:  "path",
			Usage: "stored path of file",
		},
	},
	Action: func(cctx *cli.Context) error {
		buf, err := os.ReadFile("address")
		if err != nil {
			return err
		}

		bucket := string(buf)
		object := cctx.String("object")
		path := cctx.String("path")

		client, err := lib.New()
		if err != nil {
			return err
		}

		header := make(map[string]string)

		header["test"] = "test"

		opts := miniogo.GetObjectOptions{}
		data, err := client.GetObject(cctx.Context, bucket, object, opts)
		if err != nil {
			log.Println(err)
		}

		fr, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer fr.Close()

		if _, err := io.Copy(fr, data); err != nil {
			return err
		}  

		return nil
	},
}

var ListObjectCmd = &cli.Command{
	Name:  "list",
	Usage: "list objects",
	Action: func(ctx *cli.Context) error {
		buf, err := os.ReadFile("address")
		if err != nil {
			return err
		}

		bucket := string(buf)

		client, err := lib.New()
		if err != nil {
			return err
		}

		objects := client.ListObjects(ctx.Context, bucket, miniogo.ListObjectsOptions{})
		for ob := range objects {
			fmt.Println(ob.Key)
		}

		return nil
	},
}

func fileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	_, _ = io.Copy(hash, file)

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ask4confirm() (bool, error) {
	var s string

	fmt.Printf("whether to upload(y/N): ")
	_, err := fmt.Scan(&s)
	if err != nil {
		return false, err
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true, nil
	}
	return false, nil
}
