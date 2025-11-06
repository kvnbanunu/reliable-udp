package client

import (
	"flag"
	"fmt"
	"log"
	"net"

	"reliable-udp/internal/utils"
)

type CArgs struct {
	TargetIP   string
	TargetPort uint
	Timeout    uint
	MaxRetries uint
}

type Client struct {
	Target     *net.UDPConn
	Timeout    uint
	MaxRetries uint
	BufSize    uint
	Log        string
}

func ParseArgs() *CArgs {
	args := CArgs{}
	flag.StringVar(&args.TargetIP, "target-ip", "127.0.0.1", "IP address of the server")
	flag.UintVar(&args.TargetPort, "target-port", 8080, "Port number of the server")
	flag.UintVar(&args.Timeout, "timeout", 5, "Timeout (in seconds) for waiting for acknowledgements")
	flag.UintVar(&args.MaxRetries, "max-retries", 5, "Maximum number of retries per message")
	flag.Parse()

	return &args
}

func (a *CArgs) HandleArgs() {
	if !utils.CheckIP(a.TargetIP) {
		usage("Invalid IP address")
	}

	if !utils.CheckPort(a.TargetPort) {
		usage("Invalid Port")
	}
}

func NewClient(args *CArgs, cfg *utils.Config) (*Client, error) {
	ct := Client{}

	addrStr := fmt.Sprintf("%s:%d", args.TargetIP, args.TargetPort)

	addr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		return nil, err
	}

	ct.Target, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	ct.Timeout = args.Timeout
	ct.MaxRetries = args.MaxRetries
	ct.BufSize = cfg.BufSize
	ct.Log = fmt.Sprintf("%sclient%s", cfg.LogPath, cfg.LogName)

	return &ct, nil
}

func (c *Client) Cleanup() {
	err := c.Target.Close()
	if err != nil {
		log.Fatalln("Failed to close socket:", err)
	}
}

func usage(msg string) {
	if msg != "" {
		log.Println(msg)
	}

	str := `Usage: Client [OPTIONS]
Options:
	--target-ip      IP address of the server
	--target-port    Port number of the server
	--timeout        Timeout (in seconds) for waiting for acknowledgements
	--max-retries    Maximum number of retries per message`

	fmt.Println(str)
}
