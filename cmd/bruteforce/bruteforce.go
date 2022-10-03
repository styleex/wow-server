package main

import (
	"log"
	"net"
	"time"
	"wow-server/pkg/hashcash"
)

func main() {
	pow := hashcash.NewHashCash(23)

	token := pow.NewToken(&net.IPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Zone: "",
	})
	log.Printf("token: %s (%x)", string(token), token)

	t1 := time.Now()
	nonce, err := pow.Bruteforce(token)
	if err != nil {
		log.Fatalf("cant bruteforce: %s", err)
	}
	log.Printf("nonce: %d (%s elapsed)", nonce, time.Since(t1))

	log.Println(pow.Verify(token, nonce))
}
