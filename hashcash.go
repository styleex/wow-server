package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"log"
)

type HashCash struct {
	target []byte
}

func NewHashCash(difficulty int) *HashCash {
	target := make([]byte, 256-difficulty)
	target[0] = 1
	return &HashCash{target}
}

func (h *HashCash) NewToken() []byte {
	ret := make([]byte, 32)
	if _, err := rand.Read(ret); err != nil {
		log.Panicf("Failed to generate token: %s", err)
	}

	return ret
}

func (h *HashCash) Verify(token []byte, nonce uint32) {
	nonceData := make([]byte, 32)
	binary.BigEndian.PutUint32(nonceData, nonce)

	hash := sha256.New()
	hash.Write(token)
	hash.Write(nonceData)

	resp := hash.Sum(nil)

	bytes.Compare()
}

func (h *HashCash) Bruteforce() {

}
