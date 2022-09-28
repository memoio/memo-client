package common

import (
	"github.com/memoio/memo-client/lib/types"
	"golang.org/x/xerrors"
)

var (
	ErrBadKeyType    = xerrors.New("invalid or unsupported key type")
	ErrBadPrivateKey = xerrors.New("invalid private key")
	ErrBadSign       = xerrors.New("invalid signature")
	ErrBadPublickKey = xerrors.New("invalid public key")
)

type Key interface {
	// Equals checks whether two PubKeys are the same
	Equals(Key) bool

	// Raw
	Raw() ([]byte, error)

	// Type returns the protobuf key type.
	Type() types.KeyType
}

type PrivKey interface {
	Key

	// Cryptographically sign the given bytes
	Sign([]byte) ([]byte, error)

	// Return a public key paired with this private key
	GetPublic() PubKey

	Deserialize([]byte) error
}

type PubKey interface {
	Key

	CompressedByte() ([]byte, error)
	// Verify that 'sig' is the signed hash of 'data'
	Verify(data []byte, sig []byte) (bool, error)

	Deserialize([]byte) error
}
