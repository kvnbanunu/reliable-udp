package proxy

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"reliable-udp/internal/tui"
	"reliable-udp/internal/utils"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type PArgs struct {
	ListenerAddr    string
	TargetAddr      string
	ClientDropRate  uint8
	ServerDropRate  uint8
	ClientDelayRate uint8
	ServerDelayRate uint8
	ClientDelayMin  uint
	ClientDelayMax  uint
	ServerDelayMin  uint
	ServerDelayMax  uint
}

type PRawArgs struct {
	ListenIP        string
	ListenPort      uint
	TargetIP        string
	TargetPort      uint
	ClientDropRate  uint
	ServerDropRate  uint
	ClientDelayRate uint
	ServerDelayRate uint
	ClientDelayMin  uint
	ClientDelayMax  uint
	ServerDelayMin  uint
	ServerDelayMax  uint
}

type Proxy struct {
	Listener        *net.UDPConn
	Target          *net.UDPConn
	ClientDropRate  uint8
	ServerDropRate  uint8
	ClientDelayRate uint8
	ServerDelayRate uint8
	ClientDelayMin  uint
	ClientDelayMax  uint
	ServerDelayMin  uint
	ServerDelayMax  uint

	ClientAddr *net.UDPAddr

	MaxLogs    int
	MaxDisplay int
	MsgSent    int
	MsgRecv    int
	MsgLog     []string
	Err        error

	Help           help.Model
	MsgSentDisplay progress.Model
	MsgRecvDisplay progress.Model
	Program        *tea.Program
}

func ParseArgs(cfg *utils.Config) *PRawArgs {
	args := PRawArgs{}
	var help bool

	flag.BoolVar(&help, "h", false, "Displays this help message")
	flag.StringVar(&args.ListenIP, "listen-ip", cfg.ProxyIP, "IP address to bind for client packets")
	flag.UintVar(&args.ListenPort, "listen-port", uint(cfg.ProxyPort), "UDP port to listen on for client packets")
	flag.StringVar(&args.TargetIP, "target-ip", cfg.ServerIP, "Server IP address to forward packets to")
	flag.UintVar(&args.TargetPort, "target-port", uint(cfg.ServerPort), "Server port number")
	flag.UintVar(&args.ClientDropRate, "client-drop", uint(cfg.ClientDropRate), "Drop chance (%) for packets from client")
	flag.UintVar(&args.ServerDropRate, "server-drop", uint(cfg.ServerDropRate), "Drop chance (%) for packets from server")
	flag.UintVar(&args.ClientDelayRate, "client-delay", uint(cfg.ClientDelayRate), "Delay chance (%) for packets from client")
	flag.UintVar(&args.ServerDelayRate, "server-delay", uint(cfg.ServerDelayRate), "Delay chance (%) for packets from server")
	flag.UintVar(&args.ClientDelayMin, "client-delay-time-min", cfg.ClientDelayTimeMin, "Minimum delay time (ms) for client packets")
	flag.UintVar(&args.ClientDelayMax, "client-delay-time-max", cfg.ClientDelayTimeMax, "Maximum delay time (ms) for client packets")
	flag.UintVar(&args.ServerDelayMin, "server-delay-time-min", cfg.ServerDelayTimeMin, "Minimum delay time (ms) for server packets")
	flag.UintVar(&args.ServerDelayMax, "server-delay-time-max", cfg.ServerDelayTimeMax, "Maximum delay time (ms) for server packets")
	flag.Parse()

	if help {
		usage("", nil)
	}

	return &args
}

func (a *PRawArgs) HandleArgs() *PArgs {
	if !utils.CheckIP(a.ListenIP) || !utils.CheckIP(a.TargetIP) {
		usage("Invalid IP address", nil)
	}

	proxyPort, err := utils.ToUInt16(a.ListenPort)
	if err != nil {
		usage("Invalid Port", err)
	}
	serverPort, err := utils.ToUInt16(a.TargetPort)
	if err != nil {
		usage("Invalid Port", err)
	}

	if a.ListenIP == a.TargetIP {
		if proxyPort == serverPort {
			usage("Server and Proxy are on the same host and port", nil)
		}
	}

	cDropRate, err := utils.ToUInt8(a.ClientDropRate)
	if err != nil || cDropRate > 100 {
		usage("Invalid Client drop rate", err)
	}
	sDropRate, err := utils.ToUInt8(a.ServerDropRate)
	if err != nil || sDropRate > 100 {
		usage("Invalid Server drop rate", err)
	}
	cDelayRate, err := utils.ToUInt8(a.ClientDelayRate)
	if err != nil || cDelayRate > 100 {
		usage("Invalid Client delay rate", err)
	}
	sDelayRate, err := utils.ToUInt8(a.ServerDelayRate)
	if err != nil || sDelayRate > 100 {
		usage("Invalid Server delay rate", err)
	}

	res := PArgs{}
	res.ListenerAddr = fmt.Sprintf("%s:%d", a.ListenIP, proxyPort)
	res.TargetAddr = fmt.Sprintf("%s:%d", a.TargetIP, serverPort)
	res.ClientDropRate = cDropRate
	res.ServerDropRate = sDropRate
	res.ClientDelayRate = cDelayRate
	res.ServerDelayRate = sDelayRate
	res.ClientDelayMin = a.ClientDelayMin
	res.ClientDelayMax = a.ClientDelayMax
	res.ServerDelayMin = a.ServerDelayMin
	res.ServerDelayMax = a.ServerDelayMax

	return &res
}

func NewProxy(args *PArgs, cfg *utils.Config) (*Proxy, error) {
	px := Proxy{}

	pAddr, err := net.ResolveUDPAddr("udp", args.ListenerAddr)
	if err != nil {
		return nil, err
	}

	px.Listener, err = net.ListenUDP("udp", pAddr)
	if err != nil {
		return nil, err
	}

	sAddr, err := net.ResolveUDPAddr("udp", args.TargetAddr)
	if err != nil {
		px.Listener.Close()
		return nil, err
	}

	px.Target, err = net.DialUDP("udp", nil, sAddr)
	if err != nil {
		px.Listener.Close()
		return nil, err
	}

	px.ClientDropRate = args.ClientDropRate
	px.ServerDropRate = args.ServerDropRate
	px.ClientDelayRate = args.ClientDelayRate
	px.ServerDelayRate = args.ServerDelayRate
	px.ClientDelayMin = args.ClientDelayMin
	px.ClientDelayMax = args.ClientDelayMax
	px.ServerDelayMin = args.ServerDelayMin
	px.ServerDelayMax = args.ServerDelayMax

	px.ClientAddr = nil

	px.MaxLogs = int(cfg.MaxLogs)
	px.Help = tui.NewHelpModel()
	px.MsgSentDisplay = tui.NewProgressModel()
	px.MsgRecvDisplay = tui.NewProgressModel()
	px.Program = nil
	px.Err = nil

	return &px, nil
}

func (p *Proxy) Cleanup() {
	p.Listener.Close()
	p.Target.Close()
}

func usage(msg string, err error) {
	if msg != "" {
		if err != nil {
			msg = utils.WrapErr(msg, err).Error()
		}
		log.Println(msg)
	}

	str := `Usage: Proxy [OPTIONS]
Options:
	-h                         Displays this help message
	--listen-ip                IP address to bind for client packets
	--listen-port              Port to listen on for client packets
	--target-ip                Server IP address to forward packets to
	--target-port              Server port number
	--client-drop              Drop chance (%) for packets from client
	--server-drop              Drop chance (%) for packets from server
	--client-delay             Delay chance (%) for packets from client
	--server-delay             Delay chance (%) for packets from server
	--client-delay-time-min    Minimum delay time (ms) for client packets
	--client-delay-time-max    Maximum delay time (ms) for client packets
	--server-delay-time-min    Minimum delay time (ms) for server packets
	--server-delay-time-max    Maximum delay time (ms) for server packets`

	fmt.Println(str)
	os.Exit(0)
}
