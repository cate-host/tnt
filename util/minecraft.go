package util

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type ServerStatus struct {
	Description interface{} `json:"description"`
	Players     struct {
		Max    int `json:"max"`
		Online int `json:"online"`
	} `json:"players"`
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Favicon string `json:"favicon,omitempty"`
}

func Query(ctx context.Context, address string, port uint16) (*ServerStatus, error) {
	dialer := net.Dialer{Timeout: 5 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()

	handshake := []byte{0x00}
	handshake = append(handshake, encodeVarInt(-1)...)
	handshake = append(handshake, encodeString(address)...)

	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, port)
	handshake = append(handshake, portBytes...)
	handshake = append(handshake, 0x01)

	if err := writePacket(conn, handshake); err != nil {
		return nil, fmt.Errorf("handshake: %w", err)
	}

	if err := writePacket(conn, []byte{0x00}); err != nil {
		return nil, fmt.Errorf("status request: %w", err)
	}

	data, err := readPacket(conn)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	_, n := decodeVarInt(data)
	jsonStr, _ := decodeString(data[n:])

	var status ServerStatus
	if err := json.Unmarshal([]byte(jsonStr), &status); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}
	return &status, nil
}
func encodeVarInt(val int) []byte {
	uval := uint32(val)
	var buf []byte
	for {
		b := byte(uval & 0x7F)
		uval >>= 7
		if uval != 0 {
			b |= 0x80
		}
		buf = append(buf, b)
		if uval == 0 {
			break
		}
	}
	return buf
}

func decodeVarInt(data []byte) (int, int) {
	var val uint32
	for i := 0; i < 5 && i < len(data); i++ {
		val |= uint32(data[i]&0x7F) << (7 * i)
		if data[i]&0x80 == 0 {
			return int(int32(val)), i + 1
		}
	}
	return 0, 0
}

func encodeString(s string) []byte {
	b := encodeVarInt(len(s))
	return append(b, []byte(s)...)
}

func decodeString(data []byte) (string, int) {
	length, n := decodeVarInt(data)
	return string(data[n : n+length]), n + length
}

func writePacket(conn net.Conn, payload []byte) error {
	length := encodeVarInt(len(payload))
	_, err := conn.Write(append(length, payload...))
	return err
}

func readPacket(conn net.Conn) ([]byte, error) {
	var length int
	for shift := 0; shift < 35; shift += 7 {
		b := make([]byte, 1)
		if _, err := conn.Read(b); err != nil {
			return nil, err
		}
		length |= int(b[0]&0x7F) << shift
		if b[0]&0x80 == 0 {
			break
		}
	}

	data := make([]byte, length)
	read := 0
	for read < length {
		n, err := conn.Read(data[read:])
		if err != nil {
			return nil, err
		}
		read += n
	}
	return data, nil
}
