package proxy

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"reliable-udp/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type PArgs struct {
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
	ClientDropRate  uint
	ServerDropRate  uint
	ClientDelayRate uint
	ServerDelayRate uint
	ClientDelayMin  uint
	ClientDelayMax  uint
	ServerDelayMin  uint
	ServerDelayMax  uint
	BufSize         uint
	Log             string
}

func ParseArgs() *PArgs {
	args := PArgs{}
	var help bool

	flag.BoolVar(&help, "h", false, "Displays this help message")
	flag.StringVar(&args.ListenIP, "listen-ip", "127.0.0.1", "IP address to bind for client packets")
	flag.UintVar(&args.ListenPort, "listen-port", 8081, "UDP port to listen on for client packets")
	flag.StringVar(&args.TargetIP, "target-ip", "127.0.0.1", "Server IP address to forward packets to")
	flag.UintVar(&args.TargetPort, "target-port", 8080, "Server port number")
	flag.UintVar(&args.ClientDropRate, "client-drop", 10, "Drop chance (%) for packets from client")
	flag.UintVar(&args.ServerDropRate, "server-drop", 5, "Drop chance (%) for packets from server")
	flag.UintVar(&args.ClientDelayRate, "client-delay", 20, "Delay chance (%) for packets from client")
	flag.UintVar(&args.ServerDelayRate, "server-delay", 15, "Delay chance (%) for packets from server")
	flag.UintVar(&args.ClientDelayMin, "client-delay-time-min", 100, "Minimum delay time (ms) for client packets")
	flag.UintVar(&args.ClientDelayMax, "client-delay-time-max", 200, "Maximum delay time (ms) for client packets")
	flag.UintVar(&args.ServerDelayMin, "server-delay-time-min", 150, "Minimum delay time (ms) for server packets")
	flag.UintVar(&args.ServerDelayMax, "server-delay-time-max", 300, "Maximum delay time (ms) for server packets")
	flag.Parse()

	if help {
		usage("")
	}

	return &args
}

func (a *PArgs) HandleArgs() {
	if !utils.CheckIP(a.ListenIP) || !utils.CheckIP(a.TargetIP) {
		usage("Invalid IP address")
	}

	if !utils.CheckPort(a.ListenPort) || !utils.CheckPort(a.TargetPort) {
		usage("Invalid Port")
	}
}

func NewProxy(args *PArgs, cfg *utils.Config) (*Proxy, error) {
	px := Proxy{}

	addrStr := fmt.Sprintf("%s:%d", args.ListenIP, args.ListenPort)

	addr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		return nil, err
	}

	px.Listener, err = net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	addrStr = fmt.Sprintf("%s:%d", args.TargetIP, args.TargetPort)

	addr, err = net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		px.Listener.Close()
		return nil, err
	}

	px.Target, err = net.DialUDP("udp", nil, addr)
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
	px.BufSize = cfg.BufSize
	px.Log = fmt.Sprintf("%sproxy%s", cfg.LogPath, cfg.LogName)

	return &px, nil
}

func (p *Proxy) Cleanup() {
	err := p.Listener.Close()
	if err != nil {
		log.Fatalln("Failed to close socket:", err)
	}
	err = p.Target.Close()
	if err != nil {
		log.Fatalln("Failed to close socket:", err)
	}
}

func usage(msg string) {
	if msg != "" {
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

func (p Proxy) Init() tea.Cmd {
	return nil
}

func (p Proxy) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return p, nil
}

func (p Proxy) View() string {
	var view string
	return view
}
