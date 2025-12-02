package main

import (
	"fmt"
	"log"

	"reliable-udp/internal/proxy"
	"reliable-udp/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load Config:", err)
	}

	rawArgs := proxy.ParseArgs(cfg)
	args := rawArgs.HandleArgs()

	px, err := proxy.NewProxy(args, cfg)
	if err != nil {
		log.Fatalln("Failed to setup client:", err)
	}

	defer px.Cleanup()

	m := tea.NewProgram(px)
	px.Program = m

	_, err = m.Run()
	if err != nil {
		fmt.Println("Error running proxy model:", err)
	}
}
