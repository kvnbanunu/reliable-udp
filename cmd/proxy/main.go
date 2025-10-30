package main

import (
	"log"

	"reliable-udp/internal/proxy"
	"reliable-udp/internal/utils"
)

func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load Config:", err)
	}

	args := utils.ParseArgs(utils.PROXY)

	px, err := proxy.NewProxy(args, cfg)
	if err != nil {
		log.Fatalln("Failed to setup client:", err)
	}

	defer px.Cleanup()
}
