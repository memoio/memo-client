package wallet

import (
	"context"
	"sort"
	"strings"
	"sync"

	"github.com/memoio/memo-client/lib/address"
	"github.com/memoio/memo-client/lib/crypto/signature"
	"github.com/memoio/memo-client/lib/crypto/signature/common"
	"github.com/memoio/memo-client/lib/types"
	"github.com/memoio/memo-client/lib/utils"
	"golang.org/x/xerrors"
)

type LocalWallet struct {
	lw       sync.Mutex
	password string // used for decrypt; todo plaintext is not good
	accounts map[address.Address]common.PrivKey
	keystore types.KeyStore // store
}

func New(pw string, ks types.KeyStore) *LocalWallet {
	lw := &LocalWallet{
		password: pw,
		keystore: ks,
		accounts: make(map[address.Address]common.PrivKey),
	}

	return lw
}

func (w *LocalWallet) WalletImport(ctx context.Context, ki *types.KeyInfo) (address.Address, error) {
	switch ki.Type {
	case types.Secp256k1, types.BLS:
		privkey, err := signature.ParsePrivateKey(ki.SecretKey, ki.Type)
		if err != nil {
			return address.Undef, err
		}

		pubKey := privkey.GetPublic()

		cbyte, err := pubKey.Raw()
		if err != nil {
			return address.Undef, err
		}

		addr, err := address.NewAddress(cbyte)
		if err != nil {
			return address.Undef, err
		}

		err = w.keystore.Put(addr.String(), w.password, *ki)
		if err != nil {
			return address.Undef, err
		}

		w.lw.Lock()
		w.accounts[addr] = privkey
		w.lw.Unlock()

		// for eth short addr
		if ki.Type == types.Secp256k1 {
			addrByte := utils.ToEthAddress(cbyte)

			eaddr, err := address.NewAddress(addrByte)
			if err != nil {
				return address.Undef, err
			}
			err = w.keystore.Put(eaddr.String(), w.password, *ki)
			if err != nil {
				return address.Undef, err
			}

			w.lw.Lock()
			w.accounts[eaddr] = privkey
			w.lw.Unlock()
		}

		return addr, nil
	default:
		return address.Undef, xerrors.New("unsupported key type")
	}

}

func (w *LocalWallet) WalletList(ctx context.Context) ([]address.Address, error) {
	as, err := w.keystore.List()
	if err != nil {
		return nil, err
	}

	out := make([]address.Address, 0, len(as))

	for _, s := range as {
		if strings.HasPrefix(s, address.AddrPrefix) {
			addr, err := address.NewFromString(s)
			if err != nil {
				continue
			}

			out = append(out, addr)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].String() < out[j].String()
	})

	return out, nil
}
