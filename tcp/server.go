package tcp

import (
	"log"
	"net"
	"time"
)

type Client struct {
	Conn        net.Conn
	LastMessage time.Time
}

type Server struct {
	tcpListener net.Listener

	clients map[string]*Client
}

func NewServer(ServerPort string) *Server {
	tcpListener, err := net.Listen("tcp", ServerPort)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Listening TCP Server on %s\n", ServerPort)

	return &Server{
		tcpListener: tcpListener,
		clients:     make(map[string]*Client),
	}
}

func (s *Server) RunServer() {
	for {
		conn, err := s.tcpListener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		s.acceptConnections(conn)
		go s.receiveMessages(s.clients[conn.RemoteAddr().String()])
	}
}

// 연결을 받아 처리하는 함수
func (s *Server) acceptConnections(conn net.Conn) {
	// 클라이언트 정보 추출
	clientAddr := conn.RemoteAddr().String()

	// 새로운 클라이언트 생성 및 맵에 추가
	client := &Client{
		Conn:        conn,
		LastMessage: time.Now(),
	}
	s.clients[clientAddr] = client

	log.Printf("Accepted connection from %s\n", clientAddr)
}

// 클라이언트로부터 받은 메시지를 처리하는 부분
func (s *Server) receiveMessages(client *Client) {
	// 클라이언트의 연결이 종료될 때까지 계속해서 메시지를 받습니다.
	for {
		buffer := make([]byte, 1024)
		n, err := client.Conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from %s: %v\n", client.Conn.RemoteAddr().String(), err)
			break
		}

		if n > 0 {
			message := string(buffer[:n])
			log.Printf("Received from %s: %s", client.Conn.RemoteAddr().String(), message)

			// 받은 메시지를 다른 클라이언트들에게 뿌려주기
			s.broadcast(message, client.Conn)
		}

		// 메시지를 받은 시간 업데이트
		client.LastMessage = time.Now()
	}
}

// 다른 클라이언트들에게 메시지를 뿌리는 함수
func (s *Server) broadcast(message string, sender net.Conn) {
	for _, otherClient := range s.clients {
		// 메시지를 보낸 클라이언트에게는 전송하지 않음
		if otherClient.Conn != sender {
			_, err := otherClient.Conn.Write([]byte(message))
			if err != nil {
				log.Printf("Error writing to %s: %v\n", otherClient.Conn.RemoteAddr().String(), err)
			}
		}
	}
}
