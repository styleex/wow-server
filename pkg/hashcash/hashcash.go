package hashcash

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net"
	"time"
)

type HashCash struct {
	ComplexityLevel int
	target          *big.Int
}

// NewHashCash complexityLevel in [1..256]
func NewHashCash(complexityLevel int) *HashCash {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-complexityLevel))

	return &HashCash{complexityLevel, target}
}

func encodeUInt64(val uint64) []byte {
	ret := make([]byte, 8)
	binary.BigEndian.PutUint64(ret, val)
	return ret
}

func (h *HashCash) NewToken(addr net.Addr) []byte {
	token := fmt.Sprintf("%s:%d:%d", addr.String(), time.Now().UnixMilli(), rand.Uint64())
	return []byte(token)
}

func (h *HashCash) Verify(token []byte, nonce uint64) bool {
	var hashInt big.Int
	hashInt.SetBytes(getHash(token, nonce))
	return hashInt.Cmp(h.target) == -1
}

func (h *HashCash) Bruteforce(token []byte) (uint64, error) {
	var nonce uint64 = 0
	for nonce < math.MaxUint64 {
		if h.Verify(token, nonce) {
			return nonce, nil
		}

		nonce += 1
	}

	return 0, fmt.Errorf("cant find nonce for token \"%x\" and ComplexityLevel %d", token, h.ComplexityLevel)
}

func getHash(token []byte, nonce uint64) []byte {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)

	hash := sha256.New()
	hash.Write(token)
	hash.Write(nonceBytes)
	return hash.Sum(nil)
}
