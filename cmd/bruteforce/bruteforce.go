package main

import (
	"log"
	"time"
	"wow-server/pkg"
)

func main() {
	pow := pkg.NewHashCash(23)
	token := pow.NewToken()
	log.Printf("%x", token)

	t1 := time.Now()
	nonce, err := pow.Bruteforce(token)
	if err != nil {
		log.Fatalf("cant bruteforce: %s", err)
	}
	log.Printf("nonce: %d (%s elapsed)", nonce, time.Since(t1))

	log.Println(pow.Verify(token, nonce))
}
