package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/memoio/memo-client/lib"
	miniogo "github.com/minio/minio-go/v7"
	"github.com/urfave/cli/v2"
)

var PutObjectCmd = &cli.Command{
	Name:  "put",
	Usage: "put object",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "bucket",
			Aliases: []string{"bn"},
			Usage:   "bucketName",
		},
		&cli.StringFlag{
			Name:    "object",
			Aliases: []string{"on"},
			Usage:   "objectName",
		},
		&cli.StringFlag{
			Name:  "path",
			Usage: "path of file",
		},
	},
	Action: func(cctx *cli.Context) error {
		bucket := cctx.String("bucket")
		object := cctx.String("object")
		path := cctx.String("path")

		fileinfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		r := bytes.NewBuffer(data)

		client, err := lib.New()
		if err != nil {
			return err
		}

		info, err := client.PutObject(cctx.Context, bucket, object, r, fileinfo.Size(), miniogo.PutObjectOptions{})
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
			Name:    "bucket",
			Aliases: []string{"bn"},
			Usage:   "bucketName",
		},
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
		bucket := cctx.String("bucket")
		object := cctx.String("object")
		path := cctx.String("path")

		client, err := lib.New()
		if err != nil {
			return err
		}

		data, err := client.GetObject(cctx.Context, bucket, object, miniogo.GetObjectOptions{})
		if err != nil {
			log.Println(err)
		}

		fr, err := os.Open(path)
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
