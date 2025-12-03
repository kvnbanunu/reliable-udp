package server

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
)

type SArgs struct {
	Addr string
}

type SRawArgs struct {
	ListenIP   string
	ListenPort uint
}

type Server struct {
	// Communication data
	Listener   *net.UDPConn
	ClientAddr *net.UDPAddr
	CurrentSeq uint8

	// Logging data
	MaxLogs int
	MsgSent int
	MsgRecv int
	MsgLog  []string
	Err     error

	// Render models
	Help           help.Model
	MsgSentDisplay progress.Model
	MsgRecvDisplay progress.Model
}

func ParseArgs(cfg *utils.Config) *SRawArgs {
	args := SRawArgs{}
	var help bool

	flag.BoolVar(&help, "h", false, "Displays this help message")
	flag.StringVar(&args.ListenIP, "listen-ip", cfg.ServerIP, "IP address to bind to")
	flag.UintVar(&args.ListenPort, "listen-port", uint(cfg.ServerPort), "UDP port to listen on")
	flag.Parse()

	if help {
		usage("", nil)
	}

	return &args
}

func (a *SRawArgs) HandleArgs() *SArgs {
	if !utils.CheckIP(a.ListenIP) {
		usage("Invalid IP address", nil)
	}

	port, err := utils.ToUInt16(a.ListenPort)
	if err != nil {
		usage("Invalid Port", err)
	}

	res := SArgs{}

	res.Addr = fmt.Sprintf("%s:%d", a.ListenIP, port)
	return &res
}

func NewServer(args *SArgs, cfg *utils.Config) (*Server, error) {
	srv := Server{}

	addr, err := net.ResolveUDPAddr("udp", args.Addr)
	if err != nil {
		return nil, err
	}

	srv.Listener, err = net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	srv.ClientAddr = nil
	srv.MaxLogs = int(cfg.MaxLogs)
	srv.Help = tui.NewHelpModel()
	srv.MsgSentDisplay = tui.NewProgressModel()
	srv.MsgRecvDisplay = tui.NewProgressModel()
	srv.Err = nil

	return &srv, nil
}

func (s *Server) Cleanup() {
	s.Listener.Close()
}

func usage(msg string, err error) {
	if msg != "" {
		if err != nil {
			msg = utils.WrapErr(msg, err).Error()
		}
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
