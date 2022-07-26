package keccak

import "github.com/m0ssc0de/scale.go/pkg/go-ethereum/crypto/sha3"

func Keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		_, _ = d.Write(b)
	}
	return d.Sum(nil)
}
