package client

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"reliable-udp/internal/tui"
	"reliable-udp/internal/utils"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
)

type CArgs struct {
	TargetIP   string
	TargetPort uint
	Timeout    uint
	MaxRetries uint
}

// Holds context for the client program
type Client struct {
	Target         *net.UDPConn  // Server connection
	LogDir         string        // Directory path to log file
	LogPath        string        // Full path to log file
	Timeout        time.Duration // Max time to wait for ack packets
	MaxRetries     uint          // Limit of packet resend attempts
	SeqNum         int           // Sequence number of the current message
	Max            int           // Higher number of sent/recv
	MsgSent        int           // Count of messages sent
	MsgRecv        int           // Count of messages received
	Err            error
	CurrentMsg     string
	CurrentPacket  utils.Packet
	Help           help.Model
	Input          textinput.Model
	MsgSentDisplay progress.Model
	MsgRecvDisplay progress.Model
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

	ct.LogDir = cfg.LogDir
	ct.LogPath = fmt.Sprintf("%sclient%s", cfg.LogDir, cfg.LogName)
	err = utils.PrepareLogFile(ct.LogPath)
	if err != nil {
		ct.Target.Close()
		return nil, err
	}

	ct.Timeout = time.Duration(args.Timeout) * time.Second
	ct.MaxRetries = args.MaxRetries
	ct.Help = tui.NewHelpModel()
	ct.Input = tui.NewTextInputModel()
	ct.MsgSentDisplay = tui.NewProgressModel()
	ct.MsgRecvDisplay = tui.NewProgressModel()
	ct.Err = nil

	return &ct, nil
}

func (c *Client) Cleanup() {
	c.Target.Close()
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
