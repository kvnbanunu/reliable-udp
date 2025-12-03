package client

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"reliable-udp/internal/packet"
	"reliable-udp/internal/tui"
	"reliable-udp/internal/utils"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
)

// Holds validated command line arguments
type CArgs struct {
	Target     string // Server address and port
	Timeout    uint8  // Read timeout in seconds
	MaxRetries uint8  // Max retransmission attempts
}

// Holds raw command line arguments
type CRawArgs struct {
	TargetIP   string // Server IP address
	TargetPort uint   // Server listening Port
	Timeout    uint   // Read timeout in seconds
	MaxRetries uint   // Max retransmission attempts
}

// Holds context for the client program
type Client struct {
	// Config
	Target     *net.UDPConn // Server connection
	Timeout    uint8        // Read timeout in seconds
	MaxRetries uint8        // Limit of packet resend attempts

	// Communication data
	CurrentSeq    uint8         // Sequence number of the current message
	CurrentMsg    string        // Current input message
	CurrentPacket packet.Packet // Current packet to send

	// Logging data
	MaxLogs int      // Max number of logs to show on screen
	MsgSent int      // Count of messages sent
	MsgRecv int      // Count of messages received
	MsgLog  []string // The log message that will render to screen
	Err     error

	// Render models
	Help           help.Model      // Displays controls
	Input          textinput.Model // User input
	MsgSentDisplay progress.Model  // Messages sent bar graph
	MsgRecvDisplay progress.Model  // Messages received bar graph
}

func ParseArgs(cfg *utils.Config) *CRawArgs {
	args := CRawArgs{}
	var help bool

	targetIPDefault := cfg.ServerIP
	targetPortDefault := cfg.ServerPort

	if cfg.UseProxy {
		targetIPDefault = cfg.ProxyIP
		targetPortDefault = cfg.ProxyPort
	}

	flag.BoolVar(&help, "h", false, "Displays this help message")
	flag.StringVar(&args.TargetIP, "target-ip", targetIPDefault, "IP address of the server")
	flag.UintVar(&args.TargetPort, "target-port", uint(targetPortDefault), "Port number of the server")
	flag.UintVar(&args.Timeout, "timeout", uint(cfg.Timeout), "Timeout (in seconds) for waiting for acknowledgements")
	flag.UintVar(&args.MaxRetries, "max-retries", uint(cfg.MaxRetries), "Maximum number of retries per message")
	flag.Parse()

	if help {
		usage("", nil)
	}

	return &args
}

// Converts the row argument values into proper args
func (a *CRawArgs) HandleArgs() *CArgs {
	if !utils.CheckIP(a.TargetIP) {
		usage("Invalid IP address", nil)
	}

	port, err := utils.ToUInt16(a.TargetPort)
	if err != nil {
		usage("Invalid Port", err)
	}

	res := CArgs{}

	res.Target = fmt.Sprintf("%s:%d", a.TargetIP, port)
	res.Timeout, err = utils.ToUInt8(a.Timeout)
	if err != nil {
		usage("Invalid Timeout", err)
	}

	res.MaxRetries, err = utils.ToUInt8(a.MaxRetries)
	if err != nil {
		usage("Invalid Max Retries", err)
	}

	return &res
}

// Returns Client struct with initialized fields
func NewClient(args *CArgs, cfg *utils.Config) (*Client, error) {
	ct := Client{}

	addr, err := net.ResolveUDPAddr("udp", args.Target)
	if err != nil {
		return nil, err
	}

	ct.Target, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	ct.Timeout = args.Timeout
	ct.MaxRetries = args.MaxRetries
	ct.MaxLogs = int(cfg.MaxLogs)
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

// Print usage message and exit program
func usage(msg string, err error) {
	if msg != "" {
		if err != nil {
			msg = utils.WrapErr(msg, err).Error()
		}
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
