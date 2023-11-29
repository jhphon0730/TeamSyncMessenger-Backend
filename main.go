package main

import (
	"TeamSyncMessenger-Backend/routes"
	"TeamSyncMessenger-Backend/utils"
	"log"
)

func main() {
	PORT := ":8080"
	r := routes.SetupRouter()

	log.Printf("Listening Server: %s%s\n", utils.GetOutboundIP(), PORT)
	if err := r.Run(PORT); err != nil {
		log.Panic(err.Error())
		return
	}
}
