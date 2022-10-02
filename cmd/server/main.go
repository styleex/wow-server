package main

import (
	"fmt"
	"log"
	"net"
	"wow-server/pkg"
)

type RequestMessage struct {
	// Only "GET" accepted
	Method string `json:"version"`
}

type ChallengeMessage struct {
	Token []byte `json:"token"`
}

type ChallengeResponseMessage struct {
	Nonce uint64 `json:"nonce"`
}

type PayloadResponse struct {
	Payload string `json:"payload"`
}

type Server struct {
	ddosProtector pkg.HashCash
}

func (s *Server) Listen(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panicf("Failed to listen: %s", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			// TODO: Check temporary errors?
			log.Panicf("Failed to accept connection: %s", err)
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

	var request RequestMessage
	if err := pkg.ReadData(conn, &request); err != nil {
		return err
	}

	if request.Method != "GET" {
		return fmt.Errorf("invalid method: %s", request.Method)
	}

	token := s.ddosProtector.NewToken()
	if err := pkg.WriteData(conn, ChallengeMessage{token}); err != nil {
		return err
	}

	var challengeResponse ChallengeResponseMessage
	if err := pkg.ReadData(conn, &challengeResponse); err != nil {
		return err
	}

	if !s.ddosProtector.Verify(token, challengeResponse.Nonce) {
		return fmt.Errorf("nonce not valid")
	}

	if err := pkg.WriteData(conn, PayloadResponse{""}); err != nil {
		return err
	}

	return nil
}

func main() {
	server := Server{}

	fmt.Println("Listen on :8081")
	server.Listen(":8081")
}
