package client

import (
	"fmt"

	"reliable-udp/internal/packet"
	"reliable-udp/internal/utils"
)

func (c *Client) handleInput(msg string) error {
	ptype := packet.SND
	var err error

	c.CurrentMsg = msg

	if c.CurrentSeq == 0 {
		ptype = packet.SYN
	}

	c.CurrentPacket, err = packet.NewPacket(c.CID, c.CurrentSeq, ptype, c.Timeout, msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) onSent() {
	c.addLog(fmt.Sprintf("Sent Packet: SEQ %d", c.CurrentSeq))
	c.MsgSent++
}

func (c *Client) onTimeout() {
	if c.CurrentPacket.RET == c.MaxRetries {
		c.resetState()
		c.addLog("Max retries reached. Please enter a new message")
		c.Err = packet.ErrCancel
		return
	}
	c.CurrentPacket.RET++
	c.addLog("Read timed out, resending...")
}

func (c *Client) onRecv(p packet.Packet) {
	c.MsgRecv++
	c.addLog(fmt.Sprintf("Packet Received: SEQ %d", p.SEQ))
	// the server or proxy should only ever send ACKs
	if p.TYP != packet.ACK {
		c.Err = utils.WrapErr("handleRecv", packet.ErrInvTYP)
		c.addLog(packet.ErrInvTYP.Error())
		return
	}

	if p.SEQ != c.CurrentSeq {
		c.Err = utils.WrapErr("handleRecv", packet.ErrDupPCK)
		c.addLog(packet.ErrDupPCK.Error())
		return
	}

	// connection established
	if p.SEQ == 0 {
		c.CID = p.CID
		c.addLog(fmt.Sprintf("Connection established, ID set to %d", c.CID))
	}

	c.addLog(fmt.Sprintf("Successful transmission of Packet: SEQ %d", c.CurrentSeq))
	c.addLog("Please enter a new message")
	c.resetState()
}

func (c *Client) resetState() {
	c.CurrentSeq++
	c.CurrentMsg = ""
	c.CurrentPacket = packet.Packet{}
}

// Adds log string to list and shifts the list
func (c *Client) addLog(str string) {
	if len(c.MsgLog) == c.MaxLogs {
		c.MsgLog = append(c.MsgLog[1:], str)
		return
	}
	c.MsgLog = append(c.MsgLog, str)
}
