package tui

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"reliable-udp/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func LogMessageCmd(logDir, prog, logPath string, logMsg LogMsg) tea.Cmd {
	return func() tea.Msg {
		msg, err := json.Marshal(logMsg)
		if err != nil {
			return ErrMsg{err: err}
		}

		utils.AtomicWrite(logDir, prog, logPath, msg)
		return nil
	}
}

func SendMessageCmd(conn *net.UDPConn, packet utils.Packet) tea.Cmd {
	return func() tea.Msg {
		buf := utils.Encode(packet)

		bytes, err := conn.Write(buf)
		if err != nil {
			return ErrMsg{err: err}
		}
		if bytes == 0 {
			return ErrMsg{err: fmt.Errorf("Error sending message")}
		}
		return SentMsg{}
	}
}

func RecvMessageCmd(conn *net.UDPConn) tea.Cmd {
	return func() tea.Msg {
		buf := make([]byte, utils.MAX_PACKET_LEN)

		bytes, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			return ErrMsg{err: err}
		}

		if bytes == 0 {
			return ErrMsg{err: fmt.Errorf("Error reading message")}
		}

		p := utils.Decode(buf)
		return RecvMsg{Packet: p}
	}
}

func RecvMessageTimeoutCmd(conn *net.UDPConn, timeout int) tea.Cmd {
	return func() tea.Msg {
		buf, err := utils.ReadTimeout(conn, timeout)
		if err != nil {
			if errors.Is(utils.ErrTimeout, err) {
				return TimeoutMsg{}
			}
			return ErrMsg{err: err}
		}

		p := utils.Decode(buf)
		return RecvMsg{Packet: p}
	}
}
