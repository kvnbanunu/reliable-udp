package client

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"reliable-udp/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type CArgs struct {
	TargetIP   string
	TargetPort uint
	Timeout    uint
	MaxRetries uint
}

// Holds context for the client program
type Client struct {
	Target     *net.UDPConn // Server connection
	Timeout    uint         // Max time to wait for ack packets
	MaxRetries uint         // Limit of packet resend attempts
	BufSize    uint         // Size of the buffer for read & write
	Log        string       // Path to the log file
}

func ParseArgs() *CArgs {
	args := CArgs{}
	var help bool

	flag.BoolVar(&help, "h", false, "Displays this help message")
	flag.StringVar(&args.TargetIP, "target-ip", "127.0.0.1", "IP address of the server")
	flag.UintVar(&args.TargetPort, "target-port", 8080, "Port number of the server")
	flag.UintVar(&args.Timeout, "timeout", 5, "Timeout (in seconds) for waiting for acknowledgements")
	flag.UintVar(&args.MaxRetries, "max-retries", 5, "Maximum number of retries per message")
	flag.Parse()

	if help {
		usage("")
	}

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
	-h               Displays this help message
	--target-ip      IP address of the server
	--target-port    Port number of the server
	--timeout        Timeout (in seconds) for waiting for acknowledgements
	--max-retries    Maximum number of retries per message`

	fmt.Println(str)
	os.Exit(0)
}

func (c Client) Init() tea.Cmd {
	return nil
}

func (c Client) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c Client) View() string {
	var view string
	return view
}
