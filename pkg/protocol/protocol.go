package protocol

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"time"
)

const tcpReadWriteTimeout = 30 * time.Second

type RequestMessage struct {
	// Only "GET" string accepted
	Method string `json:"method"`
}

type ChallengeMessage struct {
	Token           []byte `json:"token"`
	ComplexityLevel int    `json:"complexity_level"`
}

type ChallengeResponseMessage struct {
	Nonce uint64 `json:"nonce"`
}

type PayloadResponseMessage struct {
	Payload string `json:"payload"`
}

func WriteMessage(conn net.Conn, msg interface{}) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := conn.SetWriteDeadline(time.Now().Add(tcpReadWriteTimeout)); err != nil {
		return err
	}
	buf := make([]byte, 2+len(msgBytes))
	binary.BigEndian.PutUint16(buf, uint16(len(msgBytes)))
	copy(buf[2:], msgBytes)
	_, err = conn.Write(buf)
	return err
}

func ReadMessage(conn net.Conn, msg interface{}) error {
	if err := conn.SetReadDeadline(time.Now().Add(tcpReadWriteTimeout)); err != nil {
		return err
	}

	var buf [2]byte
	if _, err := io.ReadFull(conn, buf[:]); err != nil {
		return err
	}

	n := binary.BigEndian.Uint16(buf[:])
	msgBytes := make([]byte, n)
	if _, err := io.ReadFull(conn, msgBytes); err != nil {
		return err
	}

	return json.Unmarshal(msgBytes, msg)
}
