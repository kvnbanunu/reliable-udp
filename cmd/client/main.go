package main

import (
	"log"

	"reliable-udp/internal/client"
	"reliable-udp/internal/utils"
)

func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load Config:", err)
	}

	args := client.ParseArgs()
	args.HandleArgs()

	ct, err := client.NewClient(args, cfg)
	if err != nil {
		log.Fatalln("Failed to setup client:", err)
	}

	defer ct.Cleanup()

	// err = utils.Run(ct)
	// if err != nil {
	// 	fmt.Println("Error running client model: %v", err)
	// }
}
