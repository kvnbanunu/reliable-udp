package client

/*
This file includes all logic related to the tui for the client
*/

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
			c.Input.Blur()
			c.handleInput(c.Input.Value())
			return c, tui.SendCmd(c.Target, nil, c.CurrentPacket)
		}
	case tui.ErrMsg:
		c.Err = msg.Err
		return c, nil

	case tui.UpdateMsg:
		return c.updateProgress()
	case tui.SendSuccessMsg:
		c.onSent()
		return c, tea.Batch(tui.UpdateCmd(), tui.RecvCmd(c.Target, c.Timeout))
	case tui.TimeoutMsg:
		success := c.onTimeout()
		if !success {
			return c, tui.CancelCmd()
		}
		return c, tui.SendCmd(c.Target, nil, c.CurrentPacket)
	case tui.RecvSuccessMsg:
		success := c.onRecv(msg.Packet)
		if !success {
			return c, tea.Batch(tui.UpdateCmd(), tui.RecvCmd(c.Target, c.Timeout))
		}
		c.inputFocus()
		return c, tui.UpdateCmd()
	case tui.CancelMsg:
		c.inputFocus()
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
		strings.Join(c.MsgLog, "\n"),
		"",
		c.Help.View(tui.Keys),
	)

	return view
}

func (c Client) updateProgress() (Client, tea.Cmd) {
	c.MaxDisplay = max(c.MsgSent, c.MsgRecv)
	var cmds []tea.Cmd
	cmds = append(cmds, c.MsgSentDisplay.SetPercent(float64(c.MsgSent)/float64(c.MaxDisplay)))
	cmds = append(cmds, c.MsgRecvDisplay.SetPercent(float64(c.MsgRecv)/float64(c.MaxDisplay)))
	return c, tea.Batch(cmds...)
}

func (c *Client) inputFocus() {
	c.Input.Reset()
	c.Input.Focus()
}
