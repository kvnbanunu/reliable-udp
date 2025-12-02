package tui

import (
	"net"

	"reliable-udp/internal/packet"
)

type ErrMsg struct {
	Err error
}

type SendSuccessMsg struct{}

type RecvSuccessMsg struct {
	Packet packet.Packet
	Client *net.UDPAddr
}

type TimeoutMsg struct{}

type CancelMsg struct{}

type UpdateMsg struct{}
