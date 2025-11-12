package server

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"reliable-udp/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type SArgs struct {
	ListenIP   string
	ListenPort uint
}

type Server struct {
	Listener *net.UDPConn
	BufSize  uint
	Log      string
}

func ParseArgs() *SArgs {
	args := SArgs{}
	var help bool

	flag.BoolVar(&help, "h", false, "Displays this help message")
	flag.StringVar(&args.ListenIP, "listen-ip", "127.0.0.1", "IP address to bind to")
	flag.UintVar(&args.ListenPort, "listen-port", 8080, "UDP port to listen on")
	flag.Parse()

	if help {
		usage("")
	}

	return &args
}

func (a *SArgs) HandleArgs() {
	if !utils.CheckIP(a.ListenIP) {
		usage("Invalid IP address")
	}

	if !utils.CheckPort(a.ListenPort) {
		usage("Invalid Port")
	}
}

func NewServer(args *SArgs, cfg *utils.Config) (*Server, error) {
	srv := Server{}

	addrStr := fmt.Sprintf("%s:%d", args.ListenIP, args.ListenPort)

	addr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		return nil, err
	}

	srv.Listener, err = net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	srv.BufSize = cfg.BufSize
	srv.Log = fmt.Sprintf("%sserver%s", cfg.LogPath, cfg.LogName)

	return &srv, nil
}

func (s *Server) Cleanup() {
	err := s.Listener.Close()
	if err != nil {
		log.Fatalln("Failed to close socket:", err)
	}
}

func usage(msg string) {
	if msg != "" {
		log.Println(msg)
	}

	str := `Usage: Server [OPTIONS]
Options:
	-h               Displays this help message
	--listen-ip      IP address to bind to
	--listen-port    UDP port to listen on`

	fmt.Println(str)
	os.Exit(0)
}

func (s Server) Init() tea.Cmd {
	return nil
}

func (s Server) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s Server) View() string {
	var view string
	return view
}
