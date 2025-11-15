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

	c.resetState()
	c.MsgRecv++

	return c, nil
}

func (c Client) handleTimeout() (Client, tea.Cmd) {
	if c.CurrentPacket.Retries == uint8(c.MaxRetries) {
		return c, func() tea.Msg { return tui.CancelMsg{} }
	}

	c.CurrentPacket.Retries++
	return c, tui.SendMessageCmd(c.Target, c.CurrentPacket)
}

func (c *Client) resetState() {
	c.CurrentMsg = ""
	c.CurrentPacket = utils.Packet{}
	c.Input.Reset()
	c.Input.Focus()
	c.SeqNum++
}

func constructLog(msg tui.SendSuccessMsg) string {
	if msg.Timeout {
		return fmt.Sprintf("Sending Message (Seq #%d) retries: %d", msg.SeqNum, msg.Retries)
	}
	return fmt.Sprintf("Sending Message (Seq #%d)", msg.SeqNum)
}

func (c *Client) addLog(str string) {
	if len(c.LogMsg) == MaxLogs {
		c.LogMsg = append(c.LogMsg[1:], str)
		return
	}
	c.LogMsg = append(c.LogMsg, str)
}

func (c Client) updateProgress() (Client, tea.Cmd) {
	c.Max = max(c.MsgSent, c.MsgRecv)
	var cmds []tea.Cmd
	cmds = append(cmds, c.MsgSentDisplay.SetPercent(float64(c.MsgSent)/float64(c.Max)))
	cmds = append(cmds, c.MsgRecvDisplay.SetPercent(float64(c.MsgRecv)/float64(c.Max)))
	return c, tea.Batch(cmds...)
}
