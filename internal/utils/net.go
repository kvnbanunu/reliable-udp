package utils

import (
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	SYN uint8 = iota
	ACK
)

const (
	MAX_PAYLOAD_LEN uint8 = 255
	MAX_PACKET_LEN  int   = 258
)

const (
	SeqNumIndex int = iota
	TypeIndex
	RetriesIndex
	LengthIndex
	PayloadIndex
)

var ErrTimeout = errors.New("Error time out")

type Packet struct {
	SeqNum  uint8
	Type    uint8
	Retries uint8 // Current number of sends attempted
	Length  uint8 // Length of the payload
	Payload []byte
}

func NewPacket(seqNum, pType uint8, msg string) (Packet, error) {
	p := Packet{}
	payload := []byte(msg)
	plen := len(payload)
	if plen > int(MAX_PAYLOAD_LEN) {
		return p, fmt.Errorf("Input message too long")
	}
	p.SeqNum = seqNum
	p.Type = pType
	p.Retries = 0
	p.Length = uint8(plen)
	p.Payload = payload

	return p, nil
}

func Encode(p Packet) []byte {
	res := make([]byte, MAX_PACKET_LEN)

	res[SeqNumIndex] = p.SeqNum
	res[TypeIndex] = p.Type
	res[RetriesIndex] = p.Retries
	res[LengthIndex] = p.Length

	res = append(res[:PayloadIndex], p.Payload...)

	return res
}

func Decode(data []byte) Packet {
	res := Packet{}
	res.SeqNum = data[SeqNumIndex]
	res.Type = data[TypeIndex]
	res.Retries = data[RetriesIndex]
	res.Length = data[LengthIndex]
	res.Payload = data[PayloadIndex:]
	return res
}

func ReadTimeout(conn *net.UDPConn, timeout time.Duration) ([]byte, error) {
	buf := make([]byte, MAX_PACKET_LEN)
	err := conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, err
	}
	bytes, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, fmt.Errorf("Read failed: %w", ErrTimeout)
		}
		return nil, err
	}

	if bytes == 0 {
		return nil, fmt.Errorf("Error reading from UDP")
	}

	return buf, nil
}
