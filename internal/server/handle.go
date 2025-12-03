package server

import (
	"fmt"
	"net"

	"reliable-udp/internal/packet"
)

func (s *Server) onRecv(p packet.Packet, client *net.UDPAddr) bool {
	s.MsgRecv++
	s.addLog(fmt.Sprintf("Packet Received SEQ: %d", p.SEQ))

	switch {
	case p.TYP == packet.ACK:
		s.addLog("Bad request, client packet cannot be ACK")
		return false
	case p.SEQ < s.CurrentSeq:
		s.addLog(fmt.Sprintf("Duplicate or old message: Seq %d", p.SEQ))
		return false
	case p.SEQ == s.CurrentSeq:
		s.addLog(fmt.Sprintf("New Message Seq %d: %s", p.SEQ, p.PYL))
	case p.SEQ > s.CurrentSeq:
		s.addLog(fmt.Sprintf("Skipping Seq %d, New Message Seq %d: %s", s.CurrentSeq, p.SEQ, p.PYL))
		s.CurrentSeq = p.SEQ
	default:
		s.addLog(packet.ErrBadReq.Error())
		return false
	}

	s.ClientAddr = client
	return true
}

func (s *Server) onSent() {
	s.addLog(fmt.Sprintf("ACK sent for SEQ %d", s.CurrentSeq))
	s.MsgSent++
	s.CurrentSeq++
}

func (s *Server) addLog(str string) {
	if len(s.MsgLog) == s.MaxLogs {
		s.MsgLog = append(s.MsgLog[1:], str)
		return
	}
	s.MsgLog = append(s.MsgLog, str)
}
