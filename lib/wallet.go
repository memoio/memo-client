package lib

import (
	"sync"

	"github.com/filecoin-project/go-address"

	sig_common "github.com/memoio/memo-client/lib/crypto/signature"
	"github.com/memoio/memo-client/lib/types"
)

type LocalWallet struct {
	lw       sync.Mutex
	password string // used for decrypt; todo plaintext is not good
	accounts map[address.Address]sig_common.PrivKey
	keystore types.KeyStore // store
}

func NewW(pw string, ks types.KeyStore) *LocalWallet {
	lw := &LocalWallet{
		password: pw,
		keystore: ks,
		accounts: make(map[address.Address]sig_common.PrivKey),
	}

	return lw
}

// func WalletNew()
