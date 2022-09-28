package signature

import (
	"github.com/memoio/memo-client/lib/crypto/signature/bls"
	"github.com/memoio/memo-client/lib/crypto/signature/common"
	"github.com/memoio/memo-client/lib/crypto/signature/secp256k1"
	"github.com/memoio/memo-client/lib/types"
	"golang.org/x/xerrors"
)

func GenerateKey(typ types.KeyType) (common.PrivKey, error) {
	switch typ {
	case types.BLS:
		return bls.GenerateKey()
	case types.Secp256k1:
		return secp256k1.GenerateKey()
	default:
		return nil, common.ErrBadKeyType
	}
}

func ParsePrivateKey(privatekey []byte, typ types.KeyType) (common.PrivKey, error) {
	var privkey common.PrivKey
	switch typ {
	case types.BLS:
		privkey = &bls.PrivateKey{}
		err := privkey.Deserialize(privatekey)
		if err != nil {
			return nil, err
		}
	case types.Secp256k1:
		privkey = &secp256k1.PrivateKey{}
		err := privkey.Deserialize(privatekey)
		if err != nil {
			return nil, err
		}
	default:
		return nil, xerrors.Errorf("%d is %w", typ, common.ErrBadKeyType)
	}
	return privkey, nil
}
