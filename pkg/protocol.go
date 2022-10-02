package pkg

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"time"
)

const tcpReadWriteTimeout = 30 * time.Second

func WriteData(conn net.Conn, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := conn.SetWriteDeadline(time.Now().Add(tcpReadWriteTimeout)); err != nil {
		return err
	}
	buf := make([]byte, 2+len(dataBytes))
	binary.BigEndian.PutUint16(buf, uint16(len(dataBytes)))
	copy(buf[2:], dataBytes)
	_, err = conn.Write(buf)
	return err
}

func ReadData(conn net.Conn, out interface{}) error {
	if err := conn.SetReadDeadline(time.Now().Add(tcpReadWriteTimeout)); err != nil {
		return err
	}

	var buf [2]byte
	if _, err := io.ReadFull(conn, buf[:]); err != nil {
		return err
	}

	n := binary.BigEndian.Uint16(buf[:])
	dataBytes := make([]byte, n)
	if _, err := io.ReadFull(conn, dataBytes); err != nil {
		return err
	}

	return json.Unmarshal(dataBytes, out)
}
