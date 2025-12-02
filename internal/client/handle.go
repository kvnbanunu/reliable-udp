package client

import (
	"fmt"

	"reliable-udp/internal/packet"
)

func (c *Client) handleInput(msg string) error {
	ptype := packet.SND
	var err error

	c.CurrentMsg = msg

	c.CurrentPacket, err = packet.NewPacket(c.CurrentSeq, ptype, c.Timeout, msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) onSent() {
	c.addLog(fmt.Sprintf("Sent Packet: SEQ %d", c.CurrentSeq))
	c.MsgSent++
}

func (c *Client) onTimeout() bool {
	if c.CurrentPacket.RET == c.MaxRetries {
		c.resetState()
		c.addLog("Max retries reached. Please enter a new message")
		return false
	}
	c.CurrentPacket.RET++
	c.addLog(fmt.Sprintf("Read timed out, resending. Seq %d Retry %d", c.CurrentSeq, c.CurrentPacket.RET))
	return true
}

func (c *Client) onRecv(p packet.Packet) bool {
	c.MsgRecv++
	c.addLog(fmt.Sprintf("Packet Received: SEQ %d", p.SEQ))
	// the server or proxy should only ever send ACKs
	if p.TYP != packet.ACK {
		c.addLog("Invalid packet type, Server should only send ACK")
		return false
	}

	if p.SEQ != c.CurrentSeq {
		c.addLog(fmt.Sprintf("Duplicate or old packet received: SEQ %d", p.SEQ))
		return false
	}

	c.addLog(fmt.Sprintf("ACK received for Packet: SEQ %d", c.CurrentSeq))
	c.addLog("Please enter a new message")
	c.resetState()
	return true
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
