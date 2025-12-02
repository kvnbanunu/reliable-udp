package proxy

import (
	"net"
	"reliable-udp/internal/packet"
	"reliable-udp/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func (p *Proxy) onClientRecv(addr *net.UDPAddr) {
	p.MsgRecv++
	p.ClientAddr = addr
}

func (p *Proxy) onServerRecv() {
	p.MsgRecv++
}

func (p *Proxy) onSend() {
	p.addLog("Packet forwarded")
	p.MsgSent++
}

func ProxyListen(listener, target *net.UDPConn, p *tea.Program) tea.Cmd {
	return func() tea.Msg {
		go func() {
			for {
				packet, addr, err := packet.Recv(listener, 0)
				if err != nil {
					p.Send(tui.ErrMsg{Err: err})
					return
				}
				p.Send(recvClientMsg{Addr: addr, Packet: packet})
			}
		}()

		go func() {
			for {
				packet, _, err := packet.Recv(target, 0)
				if err != nil {
					p.Send(tui.ErrMsg{Err: err})
					return
				}
				p.Send(recvServerMsg{Packet: packet})
			}
		}()

		return nil
	}
}

type recvClientMsg struct {
	Addr *net.UDPAddr
	Packet packet.Packet
}

type recvServerMsg struct {
	Packet packet.Packet
}

func (p *Proxy) addLog(str string) {
	if len(p.MsgLog) == p.MaxLogs {
	p.MsgLog = append(p.MsgLog[1:], str)
		return
	}
	p.MsgLog = append(p.MsgLog, str)
}
