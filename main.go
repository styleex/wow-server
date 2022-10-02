package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type Server struct {
}

func (s *Server) handleConn(conn net.Conn) error {
	defer conn.Close()

	rdr := bufio.NewReader(conn)
	for {
		if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
			return err
		}

		req, err := rdr.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read line: %s", err)
		}

		switch req {
		case "GET":
			break

		case "CHALLENGE_RESPONSE":
		default:
			return fmt.Errorf("invalid line: %s", err)
		}

	}
}

func main() {
	fmt.Println("Launching server...")

	// Устанавливаем прослушивание порта
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panicf("Failed to listen: %s", err)
	}

	server := Server{}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// TODO: Check temporary errors?
			log.Panicf("Failed to accept connection: %s", err)
		}

		go server.handleConn(conn)
	}
}
