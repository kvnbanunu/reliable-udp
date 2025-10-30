package main

import (
	"log"

	"reliable-udp/internal/server"
	"reliable-udp/internal/utils"
)

func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load Config:", err)
	}

	args := utils.ParseArgs(utils.SERVER)

	srv, err := server.NewServer(args, cfg)
	if err != nil {
		log.Fatalln("Failed to setup server:", err)
	}

	defer srv.Cleanup()
}
