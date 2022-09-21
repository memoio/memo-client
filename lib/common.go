package lib

import (
	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	AccessKey string
	SecretKey string
	EndPoint  string
)

func New() (*miniogo.Client, error) {
	optionsStaticCreds := &miniogo.Options{
		Creds:        credentials.NewStaticV4(AccessKey, SecretKey, ""),
		Secure:       false,
		BucketLookup: miniogo.BucketLookupAuto,
	}

	client, err := miniogo.New(EndPoint, optionsStaticCreds)
	if err != nil {
		return nil, err
	}
	return client, nil
}
