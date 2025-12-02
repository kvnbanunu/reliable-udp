package tui

import (
	"errors"
	"net"

	"reliable-udp/internal/packet"

	tea "github.com/charmbracelet/bubbletea"
)

func ErrCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrMsg{Err: err}
	}
}

func SendCmd(conn *net.UDPConn, receiver *net.UDPAddr, p packet.Packet) tea.Cmd {
	return func() tea.Msg {
		err := packet.Send(conn, receiver, p)
		if err != nil {
			return ErrCmd(err)
		}
		return SendSuccessMsg{}
	}
}

func RecvCmd(conn *net.UDPConn, timeout uint8) tea.Cmd {
	return func() tea.Msg {
		p, client, err := packet.Recv(conn, timeout)
		if err != nil {
			if errors.Is(err, packet.ErrTimeout) {
				return TimeoutMsg{}
			}
			return ErrMsg{Err: err}
		}
		return RecvSuccessMsg{Packet: p, Client: client}
	}
}

func TimeoutCmd() tea.Cmd {
	return func() tea.Msg {
		return TimeoutMsg{}
	}
}

func CancelCmd() tea.Cmd {
	return func() tea.Msg {
		return CancelMsg{}
	}
}

func UpdateCmd() tea.Cmd {
	return func() tea.Msg {
		return UpdateMsg{}
	}
}
