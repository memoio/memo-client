package repo

import "github.com/memoio/memo-client/lib/types"

type Repo interface {
	KeyStore() types.KeyStore

	Path() (string, error)

	Close() error

	Repo() Repo 
}
