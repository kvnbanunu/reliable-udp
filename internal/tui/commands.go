package tui

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"reliable-udp/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func ErrCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrMsg{Err: err}
	}
}

func LogMessageCmd(logDir, prog, logPath string, logMsg LogMsg) tea.Cmd {
	return func() tea.Msg {
		msg, err := json.Marshal(logMsg)
		if err != nil {
			return ErrMsg{Err: err}
		}

		utils.AtomicWrite(logDir, prog, logPath, msg)
		return LogSuccessMsg{}
	}
}

func SendMessageCmd(conn *net.UDPConn, packet utils.Packet) tea.Cmd {
	return func() tea.Msg {
		buf := utils.Encode(packet)

		bytes, err := conn.Write(buf)
		if err != nil {
			return ErrMsg{Err: err}
		}
		if bytes == 0 {
			return ErrMsg{Err: fmt.Errorf("Error sending message")}
		}
		return SendSuccessMsg{
			SeqNum:  packet.SeqNum,
			Retries: packet.Retries,
			Timeout: packet.Retries != 0,
		}
	}
}

func RecvMessageCmd(conn *net.UDPConn) tea.Cmd {
	return func() tea.Msg {
		buf := make([]byte, utils.MAX_PACKET_LEN)

		bytes, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			return ErrMsg{Err: err}
		}

		if bytes == 0 {
			return ErrMsg{Err: fmt.Errorf("Error reading message")}
		}

		p := utils.Decode(buf)
		return RecvSuccessMsg{Packet: p}
	}
}

func RecvMessageTimeoutCmd(conn *net.UDPConn, timeout time.Duration) tea.Cmd {
	return func() tea.Msg {
		buf, err := utils.ReadTimeout(conn, timeout)
		if err != nil {
			if errors.Is(err, utils.ErrTimeout) {
				return TimeoutMsg{}
			}
			return ErrMsg{Err: err}
		}

		p := utils.Decode(buf)
		return RecvSuccessMsg{Packet: p}
	}
}
