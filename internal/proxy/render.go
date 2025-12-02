package proxy

import (
	"fmt"
	"strings"

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
	// var cmd tea.Cmd
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
		p.onClientRecv(msg.Addr)
		return p, tea.Batch(tui.UpdateCmd(), tui.SendCmd(p.Target, nil, msg.Packet))
	case recvServerMsg:
		p.onServerRecv()
		return p, tea.Batch(tui.UpdateCmd(), tui.SendCmd(p.Listener, p.ClientAddr, msg.Packet))
	case tui.SendSuccessMsg:
		p.onSend()
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
	p.MaxDisplay = max(p.MsgSent, p.MsgRecv)
	var cmds []tea.Cmd
	cmds = append(cmds, p.MsgSentDisplay.SetPercent(float64(p.MsgSent)/float64(p.MaxDisplay)))
	cmds = append(cmds, p.MsgRecvDisplay.SetPercent(float64(p.MsgRecv)/float64(p.MaxDisplay)))
	return p, tea.Batch(cmds...)
}
