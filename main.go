package main

import (
	"TeamSyncMessenger-Backend/routes"
	"TeamSyncMessenger-Backend/tcp"
	"TeamSyncMessenger-Backend/utils"
	"log"
)

func main() {
	PORT := ":8080"
	TCP_PORT := ":9090"
	r, db := routes.SetupRouter()

	defer db.Close()

	server := tcp.NewServer(TCP_PORT)
	go server.RunServer()

	log.Printf("Listening Server: %s%s\n", utils.GetOutboundIP(), PORT)
	log.Println("------------------------------------------------")
	if err := r.Run(PORT); err != nil {
		log.Panic(err.Error())
		return
	}
}
