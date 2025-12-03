package proxy

import (
	"fmt"
	"net"
	"strings"
	"time"

	"reliable-udp/internal/packet"
	"reliable-udp/internal/tui"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (p Proxy) Init() tea.Cmd {
	return ProxyListen(p.Listener, p.Target, p.Program)
}

func (p Proxy) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if p.Err != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, tui.Keys.Quit) {
				return p, tea.Quit
			}
		default:
			return p, nil
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, tui.Keys.Quit):
			return p, tea.Quit
		}
	case tui.ErrMsg:
		p.Err = msg.Err
		return p, nil
	case tui.UpdateMsg:
		return p.updateProgress()
	case recvClientMsg:
		p.onClientRecv(msg.Packet, msg.Addr)
		if checkRate(p.ClientDropRate) {
			p.onClientDrop(msg.Packet)
			return p, tui.UpdateCmd()
		}
		if checkRate(p.ClientDelayRate) {
			t := determineDelay(p.ClientDelayMin, p.ClientDelayMax)
			p.onClientDelay(msg.Packet, t)
			// return p, tea.Batch(tui.UpdateCmd(), delayedSendCmd(p.Target, nil, msg.Packet, t))
			return p, tea.Batch(tui.UpdateCmd(), delayMiddle(p.Target, nil, msg.Packet, t))
		}
		return p, tea.Batch(tui.UpdateCmd(), tui.SendCmd(p.Target, nil, msg.Packet))
	case recvServerMsg:
		p.onServerRecv(msg.Packet)
		if checkRate(p.ServerDropRate) {
			p.onServerDrop(msg.Packet)
			return p, tui.UpdateCmd()
		}
		if checkRate(p.ServerDelayRate) {
			t := determineDelay(p.ServerDelayMin, p.ServerDelayMax)
			p.onServerDelay(msg.Packet, t)
			// return p, tea.Batch(tui.UpdateCmd(), delayedSendCmd(p.Listener, p.ClientAddr, msg.Packet, t))
			return p, tea.Batch(tui.UpdateCmd(), delayMiddle(p.Listener, p.ClientAddr, msg.Packet, t))
		}
		return p, tea.Batch(tui.UpdateCmd(), tui.SendCmd(p.Listener, p.ClientAddr, msg.Packet))
	case delayMiddleMsg:
		return p, tea.Batch(tui.UpdateCmd(), delayedSendCmd(msg.Conn, msg.Addr, msg.Packet, msg.Delay))
	case tui.SendSuccessMsg:
		p.onSent()
		return p, tui.UpdateCmd()
	case progress.FrameMsg:
		progressModel, cmd := p.MsgSentDisplay.Update(msg)
		p.MsgSentDisplay = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
		progressModel, cmd = p.MsgRecvDisplay.Update(msg)
		p.MsgRecvDisplay = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
}

func (p Proxy) View() string {
	if p.Err != nil {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"\nError:",
			fmt.Sprintf("%v\n", p.Err),
			"Press esc or ctrl+c to exit",
		)
	}

	sentView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		p.MsgSentDisplay.View(),
		fmt.Sprintf(" %d", p.MsgSent),
	)

	recView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		p.MsgRecvDisplay.View(),
		fmt.Sprintf(" %d", p.MsgRecv),
	)

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		"\nMessages Sent:",
		sentView,
		"Messages Received:",
		recView,
		"",
		strings.Join(p.MsgLog, "\n"),
		"",
		p.Help.View(tui.Keys),
	)

	return view
}

func (p Proxy) updateProgress() (Proxy, tea.Cmd) {
	maxDisplay := max(p.MsgSent, p.MsgRecv)
	var cmds []tea.Cmd
	cmds = append(cmds, p.MsgSentDisplay.SetPercent(float64(p.MsgSent)/float64(maxDisplay)))
	cmds = append(cmds, p.MsgRecvDisplay.SetPercent(float64(p.MsgRecv)/float64(maxDisplay)))
	return p, tea.Batch(cmds...)
}

type delayMiddleMsg struct {
	Conn *net.UDPConn
	Addr *net.UDPAddr
	Packet packet.Packet
	Delay time.Duration
}

func delayMiddle(conn *net.UDPConn, addr *net.UDPAddr, p packet.Packet, d time.Duration) tea.Cmd {
	return func() tea.Msg {
		return delayMiddleMsg{
			Conn: conn,
			Addr: addr,
			Packet: p,
			Delay: d,
		}
	}
}
