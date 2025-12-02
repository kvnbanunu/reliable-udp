package packet

import (
	"net"
	"time"

	"reliable-udp/internal/utils"
)

// Read and convert buffer to Packet
func Recv(conn *net.UDPConn, timeout uint8) (Packet, *net.UDPAddr, error) {
	buf := make([]byte, MAX_PACKET_LEN)
	var p Packet

	if timeout != 0 {
		to := time.Duration(timeout) * time.Second
		err := conn.SetReadDeadline(time.Now().Add(to))
		if err != nil {
			return p, nil, err
		}
	}

	bytes, sender, err := conn.ReadFromUDP(buf)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return p, nil, utils.WrapErr("Recv", ErrTimeout)
		}
		return p, nil, utils.WrapErr("Recv", err)
	}

	if bytes == 0 {
		return p, nil, utils.WrapErr("Recv", ErrNoRead)
	}

	p = Decode(buf)

	return p, sender, nil
}

// Convert Packet to []byte and writes to target connection
func Send(conn *net.UDPConn, receiver *net.UDPAddr, p Packet) error {
	buf := Encode(p)
	var bytes int
	var err error

	if receiver != nil {
		bytes, err = conn.WriteToUDP(buf, receiver)
	} else {
		bytes, err = conn.Write(buf)
	}

	if err != nil {
		return utils.WrapErr("Send", err)
	}

	if bytes == 0 {
		return utils.WrapErr("Send", ErrNoWrite)
	}

	return nil
}
