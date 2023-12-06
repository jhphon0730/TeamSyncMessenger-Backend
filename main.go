package main

import (
	"TeamSyncMessenger-Backend/routes"
	"TeamSyncMessenger-Backend/tcp"
	"TeamSyncMessenger-Backend/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	PORT := ":8080"
	TCP_PORT := ":9090"
	r, db := routes.SetupRouter()

	defer db.Close()

	server := tcp.NewServer(TCP_PORT)

	go server.RunServer()
	go r.Run(PORT)

	log.Printf("Listening Server: %s%s\n", utils.GetOutboundIP(), PORT)
	log.Println("------------------------------------------------")

	// 시그널 처리 함수 등록
	signalHandler := make(chan os.Signal, 1)
	signal.Notify(signalHandler, os.Interrupt, syscall.SIGTERM)

	// 시그널을 수신할 때까지 대기
	<-signalHandler

	log.Println("Received termination signal. Closing server...")

	// 서버 종료 메서드 호출
	server.CloseServer()

	log.Println("Server closed.")
}
