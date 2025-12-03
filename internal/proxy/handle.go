package proxy

import (
	"fmt"
	"math/rand/v2"
	"net"
	"time"

	"reliable-udp/internal/packet"
	"reliable-udp/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func (p *Proxy) onClientRecv(packet packet.Packet, addr *net.UDPAddr) {
	p.MsgRecv++
	p.ClientAddr = addr
	p.addLog(fmt.Sprintf("Received from Client SEQ %d", packet.SEQ))
}

func (p *Proxy) onServerRecv(packet packet.Packet) {
	p.MsgRecv++
	p.addLog(fmt.Sprintf("Received from Server SEQ %d", packet.SEQ))
}

func (p *Proxy) onClientDrop(packet packet.Packet) {
	p.addLog(fmt.Sprintf("Client Packet dropped SEQ %d", packet.SEQ))
}

func (p *Proxy) onServerDrop(packet packet.Packet) {
	p.addLog(fmt.Sprintf("Server Packet dropped SEQ %d", packet.SEQ))
}

func (p *Proxy) onClientDelay(packet packet.Packet, t time.Duration) {
	p.addLog(fmt.Sprintf("Client Packet SEQ %d delayed for %v", packet.SEQ, t))
}

func (p *Proxy) onServerDelay(packet packet.Packet, t time.Duration) {
	p.addLog(fmt.Sprintf("Server Packet SEQ %d delayed for %v", packet.SEQ, t))
}

func (p *Proxy) onSent() {
	p.addLog("Packet forwarded")
	p.MsgSent++
}

func checkRate(rate uint8) bool {
	check := rand.UintN(101)
	return check <= uint(rate)
}

func determineDelay(dmin, dmax uint) time.Duration {
	delay := rand.UintN(dmax+1) + dmin
	return time.Millisecond * time.Duration(delay)
}

func delayedSendCmd(conn *net.UDPConn, addr *net.UDPAddr, packet packet.Packet, delay time.Duration) tea.Cmd {
	time.Sleep(delay)
	return tui.SendCmd(conn, addr, packet)
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
	Addr   *net.UDPAddr
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
