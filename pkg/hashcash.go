package pkg

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
)

type HashCash struct {
	complexityLevel int
	target          *big.Int
}

// NewHashCash complexityLevel in [1..256]
func NewHashCash(complexityLevel int) *HashCash {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-complexityLevel))
	log.Printf("target: %x", target.Bytes())

	return &HashCash{complexityLevel, target}
}

func (h *HashCash) NewToken() []byte {
	ret := make([]byte, 64)
	if _, err := rand.Read(ret); err != nil {
		log.Panicf("Failed to generate token: %s", err)
	}

	return ret
}

func (h *HashCash) Verify(token []byte, nonce uint64) bool {
	var hashInt big.Int
	hashInt.SetBytes(getHash(token, nonce))
	return hashInt.Cmp(h.target) == -1
}

func (h *HashCash) Bruteforce(token []byte) (uint64, error) {
	var nonce uint64 = 0
	for nonce < math.MaxInt64 {
		if h.Verify(token, nonce) {
			return nonce, nil
		}

		nonce += 1
	}

	return 0, fmt.Errorf("cant find nonce for token \"%x\" and complexityLevel %d", token, h.complexityLevel)
}

func getHash(token []byte, nonce uint64) []byte {
	nonceBytes := make([]byte, 32)
	binary.BigEndian.PutUint64(nonceBytes, nonce)

	hash := sha256.New()
	hash.Write(token)
	hash.Write(nonceBytes)
	return hash.Sum(nil)
}
