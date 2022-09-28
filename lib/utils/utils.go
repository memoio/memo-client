package utils

import "golang.org/x/crypto/sha3"

func ToEthAddress(pubkey []byte) []byte {
	if len(pubkey) == 65 {
		d := sha3.NewLegacyKeccak256()
		d.Write(pubkey[1:])
		payload := d.Sum(nil)
		return payload[12:]
	}

	return pubkey
}
