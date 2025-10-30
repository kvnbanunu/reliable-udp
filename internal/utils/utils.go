package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

type ProgramType int

const (
	PROXY ProgramType = iota
	SERVER
	CLIENT
)

type Args struct {
	ServerIP   string
	ServerPort int
	ProxyIP    string
	ProxyPort  int
	Help       bool
}

type Config struct {
	BufSize int    `json:"bufsize"`
	LogPath string `json:"logpath"`
	LogName string `json:"logname"`
	Timeout int    `json:"timeout"`
}

func LoadConfig() (*Config, error) {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ParseArgs(pt ProgramType) *Args {
	prog := os.Args[0]

	args := Args{}
	flag.StringVar(&args.ServerIP, "i", "127.0.0.1", "IP address of the host server")
	flag.IntVar(&args.ServerPort, "p", 8080, "Port of the host server")
	if pt == PROXY {
		flag.StringVar(&args.ProxyIP, "I", "127.0.0.1", "IP address of the proxy server")
		flag.IntVar(&args.ProxyPort, "P", 8081, "Port of the proxy server")
	}
	flag.BoolVar(&args.Help, "h", false, "Prints a help message")
	flag.Parse()

	if args.Help {
		usage(prog, "", pt)
	}

	if pt == PROXY {
		// if same machine and using the same port
		if args.ServerIP == args.ProxyIP &&
			args.ServerPort == args.ProxyPort {
			usage(prog, "", pt)
		}
	}

	return &args
}

func usage(prog, msg string, pt ProgramType) {
	if msg != "" {
		log.Println(msg)
	}

	str := `Usage: %s [OPTIONS]
Options:
	-h    Display this help message
	-i    IP address of the host server
	-p    Port of the host server%s
`
	if pt == PROXY {
		fmt.Printf(str, prog, `
	-I    IP address of the proxy server
	-P    Port of the proxy server
`)
	} else {
		fmt.Printf(str, prog, "")
	}
	os.Exit(0)
}
