package client

import (
	"fmt"
	"strings"

	"reliable-udp/internal/tui"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (c Client) Init() tea.Cmd {
	return textinput.Blink
}

func (c Client) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	if c.Err != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, tui.Keys.Quit) {
				return c, tea.Quit
			}
		default:
			return c, nil
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, tui.Keys.Quit):
			return c, tea.Quit
		case key.Matches(msg, tui.Keys.Enter):
			return c.handleEnter()
		}
	case tui.ErrMsg:
		c.Err = msg.Err
		return c, nil

	case tui.SendSuccessMsg:
		c.MsgSent++
		c.addLog(constructLog(msg))
		return c, tea.Batch(c.logCmd(), tui.RecvMessageTimeoutCmd(c.Target, c.Timeout))
	case tui.LogSuccessMsg:
		return c.updateProgress()
	case tui.RecvSuccessMsg:
		c.addLog("ACK received!")
		c, cmd = c.handleRecvMsg(msg)
		return c, tea.Batch(c.logCmd(), cmd)
	case tui.TimeoutMsg:
		c.addLog("Timed out waiting for ACK. Resending...")
		return c.handleTimeout()
	case tui.CancelMsg:
		c.addLog("Max number of retries attempted. Request cancelled")
		c.resetState()
		return c, nil

	case progress.FrameMsg:
		progressModel, cmd := c.MsgSentDisplay.Update(msg)
		c.MsgSentDisplay = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
		progressModel, cmd = c.MsgRecvDisplay.Update(msg)
		c.MsgRecvDisplay = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
	}

	c.Input, cmd = c.Input.Update(msg)
	cmds = append(cmds, cmd)
	return c, tea.Batch(cmds...)
}

func (c Client) View() string {
	if c.Err != nil {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"\nError:",
			fmt.Sprintf("%v\n", c.Err),
			"Press esc or ctrl+c to exit",
		)
	}

	sentView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		c.MsgSentDisplay.View(),
		fmt.Sprintf(" %d", c.MsgSent),
	)

	recView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		c.MsgRecvDisplay.View(),
		fmt.Sprintf(" %d", c.MsgRecv),
	)

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		"\nSend a message:",
		c.Input.View(),
		"\nMessages Sent:",
		sentView,
		"Messages Received:",
		recView,
		"",
		strings.Join(c.LogMsg, "\n"),
		"",
		c.Help.View(tui.Keys),
	)

	return view
}
