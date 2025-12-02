package server

import (
	"fmt"
	"net"

	"reliable-udp/internal/packet"
)

func (s *Server) onRecv(p packet.Packet, client *net.UDPAddr) uint8 {
	s.MsgRecv++
	s.addLog(fmt.Sprintf("Packet Received CID: %d SEQ: %d", p.CID, p.SEQ))

	if (p.CID != 0 && p.SEQ == 0 && p.TYP == packet.SYN) || (p.TYP == packet.ACK) {
		s.addLog(packet.ErrBadReq.Error())
		return 0
	}

	// New client
	cid := p.CID
	if p.CID == 0 && p.SEQ == 0 && p.TYP == packet.SYN {
		s.NumClients++
		cid = uint8(s.NumClients)
		s.Clients = append(s.Clients, ClientData{
			CID:        cid,
			Addr:       client,
			CurrentSeq: 0,
		})
		s.addLog(fmt.Sprintf("New client added CID: %d", cid))
	} else {
		s.Clients[p.CID].CurrentSeq = p.SEQ
	}

	s.addLog(fmt.Sprintf("Message: %s", string(p.PYL)))
	return cid
}

func (s *Server) onSend() {
	s.addLog("ACK sent")
	s.MsgSent++
}

func (s *Server) addLog(str string) {
	if len(s.MsgLog) == s.MaxLogs {
		s.MsgLog = append(s.MsgLog[1:], str)
		return
	}
	s.MsgLog = append(s.MsgLog, str)
}
