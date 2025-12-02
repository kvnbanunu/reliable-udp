package server

import (
	"fmt"
	"strings"

	"reliable-udp/internal/packet"
	"reliable-udp/internal/tui"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (s Server) Init() tea.Cmd {
	return func() tea.Msg { return initMsg{} }
}

func (s Server) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if s.Err != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, tui.Keys.Quit) {
				return s, tea.Quit
			}
		default:
			return s, nil
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, tui.Keys.Quit):
			return s, tea.Quit
		}
	case tui.ErrMsg:
		s.Err = msg.Err
		return s, nil
	case tui.UpdateMsg:
		return s.updateProgress()
	case initMsg:
		return s, tui.RecvCmd(s.Listener, 0)
	case tui.RecvSuccessMsg:
		success := s.onRecv(msg.Packet, msg.Client)
		if !success {
			return s, tea.Batch(tui.UpdateCmd(), tui.RecvCmd(s.Listener, 0))
		}
		p, err := packet.NewPacket(s.CurrentSeq, packet.ACK, 0, "")
		if err != nil {
			s.Err = err
			return s, nil
		}
		return s, tea.Batch(tui.UpdateCmd(), tui.SendCmd(s.Listener, s.ClientAddr, p))

	case tui.SendSuccessMsg:
		s.onSend()
		return s, tea.Batch(tui.UpdateCmd(), tui.RecvCmd(s.Listener, 0))

	case progress.FrameMsg:
		progressModel, cmd := s.MsgSentDisplay.Update(msg)
		s.MsgSentDisplay = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
		progressModel, cmd = s.MsgRecvDisplay.Update(msg)
		s.MsgRecvDisplay = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

func (s Server) View() string {
	if s.Err != nil {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"\nError:",
			fmt.Sprintf("%v\n", s.Err),
			"Press esc or ctrl+c to exit",
		)
	}

	sentView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		s.MsgSentDisplay.View(),
		fmt.Sprintf(" %d", s.MsgSent),
	)

	recView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		s.MsgRecvDisplay.View(),
		fmt.Sprintf(" %d", s.MsgRecv),
	)

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		"\nMessages Sent:",
		sentView,
		"Messages Received:",
		recView,
		"",
		strings.Join(s.MsgLog, "\n"),
		"",
		s.Help.View(tui.Keys),
	)

	return view
}

func (s Server) updateProgress() (Server, tea.Cmd) {
	s.MaxDisplay = max(s.MsgSent, s.MsgRecv)
	var cmds []tea.Cmd
	cmds = append(cmds, s.MsgSentDisplay.SetPercent(float64(s.MsgSent)/float64(s.MaxDisplay)))
	cmds = append(cmds, s.MsgRecvDisplay.SetPercent(float64(s.MsgRecv)/float64(s.MaxDisplay)))
	return s, tea.Batch(cmds...)
}

type initMsg struct{}
