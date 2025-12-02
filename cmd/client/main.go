package main

import (
	"fmt"
	"log"

	"reliable-udp/internal/client"
	"reliable-udp/internal/tui"
	"reliable-udp/internal/utils"
)

func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load Config:", err)
	}

	rawArgs := client.ParseArgs(cfg)
	args := rawArgs.HandleArgs()

	ct, err := client.NewClient(args, cfg)
	if err != nil {
		log.Fatalln("Failed to setup client:", err)
	}

	defer ct.Cleanup()

	err = tui.Run(ct)
	if err != nil {
		fmt.Println("Error running client model:", err)
	}
}
