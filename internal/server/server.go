package server

import (
	"fmt"
	"log"
	"net"

	"reliable-udp/internal/utils"
)

type Server struct {
	Sock    *net.UDPConn
	BufSize int
	Log     string
	Timeout int
}

func NewServer(args *utils.Args, cfg *utils.Config) (*Server, error) {
	srv := Server{}

	addrStr := fmt.Sprintf("%s:%d", args.ServerIP, args.ServerPort)

	addr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		return nil, err
	}

	srv.Sock, err = net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	srv.BufSize = cfg.BufSize
	srv.Log = fmt.Sprintf("%sserver%s", cfg.LogPath, cfg.LogName)
	srv.Timeout = cfg.Timeout

	return &srv, nil
}

func (s *Server) Cleanup() {
	err := s.Sock.Close()
	if err != nil {
		log.Fatalln("Failed to close socket:", err)
	}
}
