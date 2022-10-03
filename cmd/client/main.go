package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
	"wow-server/pkg/hashcash"
	"wow-server/pkg/protocol"
)

func doRequest(conn net.Conn) error {
	// 1. Write request
	if err := protocol.WriteMessage(conn, protocol.RequestMessage{Method: "GET"}); err != nil {
		return err
	}

	// 2. Read challenge
	var challengeMsg protocol.ChallengeMessage
	if err := protocol.ReadMessage(conn, &challengeMsg); err != nil {
		return err
	}

	// 3. Solve challenge and send result
	hashCash := hashcash.NewHashCash(challengeMsg.ComplexityLevel)
	t1 := time.Now()
	nonce, err := hashCash.Bruteforce(challengeMsg.Token)
	if err != nil {
		return err
	}
	log.Printf("Challenge solved in %s", time.Since(t1))

	if err := protocol.WriteMessage(conn, protocol.ChallengeResponseMessage{Nonce: nonce}); err != nil {
		return err
	}

	// 4. Read payload
	var payload protocol.PayloadResponseMessage
	if err := protocol.ReadMessage(conn, &payload); err != nil {
		return err
	}

	fmt.Printf("Payload: %s\n", payload.Payload)
	return nil
}

func main() {
	serverAddr := flag.String("server", "127.0.0.1:8081", "wow-server addr")
	flag.Parse()

	conn, err := net.Dial("tcp", *serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	if err := doRequest(conn); err != nil {
		log.Fatalf("%s", err)
	}
}
