package proxy

import (
	"fmt"
	"log"
	"net"

	"reliable-udp/internal/utils"
)

type Proxy struct {
	Sock      *net.UDPConn
	Server    *net.UDPConn
	BufSize   int
	Log       string
	Timeout   int
	DelayRate int
	DropRate  int
}

func NewProxy(args *utils.Args, cfg *utils.Config) (*Proxy, error) {
	px := Proxy{}

	addrStr := fmt.Sprintf("%s:%d", args.ProxyIP, args.ProxyPort)

	addr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		return nil, err
	}

	px.Sock, err = net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	addrStr = fmt.Sprintf("%s:%d", args.ServerIP, args.ServerPort)

	addr, err = net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		px.Sock.Close()
		return nil, err
	}

	px.Server, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		px.Sock.Close()
		return nil, err
	}

	px.BufSize = cfg.BufSize
	px.Log = fmt.Sprintf("%sproxy%s", cfg.LogPath, cfg.LogName)
	px.Timeout = cfg.Timeout

	return &px, nil
}

func (p *Proxy) Cleanup() {
	err := p.Sock.Close()
	if err != nil {
		log.Fatalln("Failed to close socket:", err)
	}
}
