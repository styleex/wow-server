package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
	"wow-server/pkg/cites_storage"
	"wow-server/pkg/hashcash"
	"wow-server/pkg/protocol"
)

type Server struct {
	citesStorage  *cites_storage.Storage
	ddosProtector *hashcash.HashCash
}

func NewServer(ddosComplexityLevel int) *Server {
	return &Server{
		citesStorage:  cites_storage.Load(),
		ddosProtector: hashcash.NewHashCash(ddosComplexityLevel),
	}
}

func (s *Server) ListenAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			if err := s.handleConn(conn); err != nil {
				log.Printf("WARN! Failed to process client %s: %s", conn.RemoteAddr(), err)
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	defer conn.Close()

	// 1. Read and verify request
	var request protocol.RequestMessage
	if err := protocol.ReadMessage(conn, &request); err != nil {
		return err
	}

	if request.Method != "GET" {
		return fmt.Errorf("invalid method: %s", request.Method)
	}

	// 2. Send challenge
	token := s.ddosProtector.NewToken(conn.RemoteAddr())
	if err := protocol.WriteMessage(conn, protocol.ChallengeMessage{Token: token, ComplexityLevel: s.ddosProtector.ComplexityLevel}); err != nil {
		return err
	}

	// 3. Read and verify challenge response
	var challengeResponse protocol.ChallengeResponseMessage
	if err := protocol.ReadMessage(conn, &challengeResponse); err != nil {
		return err
	}

	if !s.ddosProtector.Verify(token, challengeResponse.Nonce) {
		return fmt.Errorf("nonce not valid")
	}

	// 4. Write payload
	if err := protocol.WriteMessage(conn, protocol.PayloadResponseMessage{Payload: s.citesStorage.RandomCite()}); err != nil {
		return err
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixMilli())
	server := NewServer(22)

	fmt.Println("Listen on :8081")
	if err := server.ListenAndServe(":8081"); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
