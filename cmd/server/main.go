package main

import (
	"fmt"
	"log"

	"reliable-udp/internal/server"
	"reliable-udp/internal/tui"
	"reliable-udp/internal/utils"
)

func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load Config:", err)
	}

	rawArgs := server.ParseArgs(cfg)
	args := rawArgs.HandleArgs()

	srv, err := server.NewServer(args, cfg)
	if err != nil {
		log.Fatalln("Failed to setup server:", err)
	}

	defer srv.Cleanup()

	err = tui.Run(srv)
	if err != nil {
		fmt.Println("Error running server model:", err)
	}

	// for i := range(10){
	// 	buf := make([]byte, 300)
	// 	bytes, client, err := srv.Listener.ReadFromUDP(buf)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	//
	// 	if bytes == 0 {
	// 		fmt.Println("No bytes read")
	// 		return
	// 	}
	//
	// 	fmt.Printf("Received:\n%d\n%d\n%d\n%d\n%s\n", buf[0], buf[1], buf[2], buf[3], buf[4:])
	//
	// 	if i % 3 == 0 {
	// 		bytes, err = srv.Listener.WriteToUDP(buf, client)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 		}
	//
	// 		if bytes == 0 {
	// 			fmt.Println("No bytes written")
	// 			return
	// 		}
	// 	}
	// }
}
