package address

import (
	"bytes"

	b58 "github.com/mr-tron/base58/base58"
	"github.com/zeebo/blake3"
	"golang.org/x/xerrors"
)

var UndefAddressString = "<empty>"

const ChecksumHashLength = 4

const MaxAddressStringLength = AddrPrefixLen + 94

type Address struct{ str string }

var Undef = Address{}

const AddrPrefix = "Me"
const AddrPrefixLen = 2

func NewAddress(payload []byte) (Address, error) {
	return newAddress(payload)
}

func newAddress(payload []byte) (Address, error) {
	buf := make([]byte, len(payload))
	copy(buf[:], payload)
	return Address{string(buf)}, nil
}

func NewFromString(s string) (Address, error) {
	return decode(s)
}

func decode(a string) (Address, error) {
	if len(a) == 0 {
		return Undef, xerrors.New("invalid address length")
	}
	if a == UndefAddressString {
		return Undef, xerrors.New("invalid address length")
	}
	if len(a) > MaxAddressStringLength || len(a) < 3 {
		return Undef, xerrors.New("invalid address length")
	}

	if string(a[0:AddrPrefixLen]) != AddrPrefix {
		return Undef, xerrors.New("unknown address type")
	}

	raw := a[AddrPrefixLen:]

	payloadcksm, err := b58.Decode(raw)
	if err != nil {
		return Undef, err
	}

	if len(payloadcksm)-ChecksumHashLength < 0 {
		return Undef, xerrors.New("invalid address checksum")
	}

	payload := payloadcksm[:len(payloadcksm)-ChecksumHashLength]
	cksm := payloadcksm[len(payloadcksm)-ChecksumHashLength:]

	if !ValidateChecksum(payload, cksm) {
		return Undef, xerrors.New("invalid address checksum")
	}

	return newAddress(payload)
}

func encode(addr Address) (string, error) {
	if addr == Undef {
		return UndefAddressString, nil
	}

	cksm := Checksum(addr.Bytes())
	strAddr := AddrPrefix + b58.Encode(append(addr.Bytes(), cksm[:]...))

	return strAddr, nil
}

func Checksum(ingest []byte) []byte {
	h := blake3.New()
	h.Write(ingest)
	res := h.Sum(nil)
	return res[:ChecksumHashLength]
}

func ValidateChecksum(ingest, expect []byte) bool {
	digest := Checksum(ingest)
	return bytes.Equal(digest, expect)
}

func (a Address) Len() int {
	return len(a.str)
}

func (a Address) Bytes() []byte {
	return []byte(a.str)
}

func (a Address) String() string {
	str, err := encode(a)
	if err != nil {
		panic(err)
	}
	return str
}
