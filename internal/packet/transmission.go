package packet

import (
	"net"
	"time"

	"reliable-udp/internal/utils"
)

func SetTimeout(conn *net.UDPConn, timeout time.Duration) error {
	err := conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}

	return nil
}

// Read and convert buffer to Packet
func Recv(receiver *net.UDPConn) (Packet, *net.UDPAddr, error) {
	buf := make([]byte, MAX_PACKET_LEN)
	var p Packet

	bytes, sender, err := receiver.ReadFromUDP(buf)
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
func Send(sender *net.UDPConn, receiver *net.UDPAddr, p Packet) error {
	buf := Encode(p)
	
	bytes, err := sender.WriteToUDP(buf, receiver)
	if err != nil {
		return utils.WrapErr("Send", err)
	}

	if bytes == 0 {
		return utils.WrapErr("Send", ErrNoWrite)
	}

	return nil
}
