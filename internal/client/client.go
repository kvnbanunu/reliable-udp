package client

import (
	"fmt"
	"log"
	"net"

	"reliable-udp/internal/utils"
)

type Client struct {
	Sock    *net.UDPConn
	BufSize int
	Log     string
	Timeout int
}

func NewClient(args *utils.Args, cfg *utils.Config) (*Client, error) {
	ct := Client{}

	addrStr := fmt.Sprintf("%s:%d", args.ServerIP, args.ServerPort)

	addr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		return nil, err
	}

	ct.Sock, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	ct.BufSize = cfg.BufSize
	ct.Log = fmt.Sprintf("%sclient%s", cfg.LogPath, cfg.LogName)
	ct.Timeout = cfg.Timeout

	return &ct, nil
}

func (c *Client) Cleanup() {
	err := c.Sock.Close()
	if err != nil {
		log.Fatalln("Failed to close socket:", err)
	}
}
