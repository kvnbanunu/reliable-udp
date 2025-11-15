package client

import (
	"fmt"

	"reliable-udp/internal/tui"
	"reliable-udp/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func (c Client) handleEnter() (Client, tea.Cmd) {
	c.CurrentMsg = c.Input.Value()
	c.Input.Blur()
	c.CurrentPacket, c.Err = utils.NewPacket(uint8(c.SeqNum), utils.SYN, c.CurrentMsg)
	if c.Err != nil {
		return c, tui.ErrCmd(c.Err)
	}
	return c, tui.SendMessageCmd(c.Target, c.CurrentPacket)
}

func (c Client) logCmd() tea.Cmd {
	return tui.LogMessageCmd(c.LogDir, "client", c.LogPath, tui.LogMsg{
		MsgSent: c.MsgSent,
		MsgRecv: c.MsgRecv,
	})
}

func (c Client) handleRecvMsg(msg tui.RecvSuccessMsg) (Client, tea.Cmd) {
	p := msg.Packet

	if c.CurrentPacket.SeqNum != p.SeqNum {
		return c, tui.ErrCmd(fmt.Errorf("Mismatch sequence number"))
	}

	c.CurrentMsg = ""
	c.Input.Reset()
	c.Input.Focus()
	c.SeqNum++
	c.MsgRecv++

	return c, nil
}

func (c Client) updateProgress() (Client, tea.Cmd) {
	c.Max = max(c.MsgSent, c.MsgRecv)
	var cmds []tea.Cmd
	cmds = append(cmds, c.MsgSentDisplay.SetPercent(float64(c.MsgSent / c.Max)))
	cmds = append(cmds, c.MsgRecvDisplay.SetPercent(float64(c.MsgRecv / c.Max)))
	return c, tea.Batch(cmds...)
}
