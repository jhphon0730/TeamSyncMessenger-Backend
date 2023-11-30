package tcp

import (
	messagetype "TeamSyncMessenger-Backend/tcp/messageType"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"time"
)

type Client struct {
	Conn        net.Conn
	LastMessage time.Time
}

func (c *Client) sendMessage(message messagetype.Message) error {
	writeData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = c.Conn.Write(writeData)
	return err
}

// 추가: 연결이 이미 닫혀 있는지 확인하는 함수
func (c *Client) isClosed() bool {
	// 추가: 연결이 이미 닫혔는지 확인
	_, err := c.Conn.Read([]byte{})
	return err != nil
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

	log.Println("------------------------------------------------")
	log.Printf("Listening TCP Server on %s\n", ServerPort)

	return &Server{
		tcpListener: tcpListener,
		clients:     make(map[string]*Client),
	}
}

func (s *Server) server(messages chan messagetype.Message) {
	for {
		msg := <-messages
		switch msg.Type {
		case "disconnect":
			clientAddr := msg.Content.(string)
			if client, ok := s.clients[clientAddr]; ok {
				// 클라이언트 연결 해제 로직 추가
				delete(s.clients, clientAddr)
				client.Conn.Close()
				log.Printf("Disconnect %s\n", clientAddr)
			}
		}
	}
}

func (s *Server) RunServer() {
	messages := make(chan messagetype.Message)
	go s.server(messages)

	for {
		conn, err := s.tcpListener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		client := s.acceptConnections(conn, messages)
		go s.receiveMessages(client, messages)
	}
}

// 연결을 받아 처리하는 함수
func (s *Server) acceptConnections(conn net.Conn, messages chan messagetype.Message) *Client {
	// 클라이언트 정보 추출
	clientAddr := conn.RemoteAddr().String()

	// 새로운 클라이언트 생성 및 맵에 추가
	client := &Client{
		Conn:        conn,
		LastMessage: time.Now(),
	}
	s.clients[clientAddr] = client

	log.Println("Connected: ", clientAddr)

	message := messagetype.Message{
		Type:    "connect_success",
		Content: clientAddr,
	}

	err := client.sendMessage(message)
	if err != nil {
		log.Panic("클라이언트 연결 에러: ", err.Error())
	}

	return client
}

// 클라이언트로부터 받은 메시지를 처리하는 부분
func (s *Server) receiveMessages(client *Client, messages chan messagetype.Message) {
	defer func() {
		// 클라이언트 연결이 종료되면 해당 클라이언트를 서버에서 제거
		clientAddr := client.Conn.RemoteAddr().String()
		delete(s.clients, clientAddr)
		log.Printf("Connection from %s closed\n", clientAddr)
	}()

	for {
		// 추가: 연결이 이미 닫혀 있는지 확인
		if client.isClosed() {
			return
		}

		// 메시지를 받은 시간 업데이트
		client.LastMessage = time.Now()

		var buffer []byte
		tempBuffer := make([]byte, 1024) // Temporary buffer

		n, err := client.Conn.Read(tempBuffer)
		if err != nil && client.Conn != nil {
			if err != io.EOF {
				log.Println("메시지 받기 실패 에러: ", err.Error())
			}
			return
		}

		buffer = append(buffer, tempBuffer[:n]...)

		// Attempt to decode the received data as JSON
		var message messagetype.Message
		decoder := json.NewDecoder(bytes.NewReader(buffer))
		if err := decoder.Decode(&message); err == nil {
			// Successfully decoded a JSON object
			s.handleMessage(buffer[:decoder.InputOffset()], len(buffer[:decoder.InputOffset()]), messages)
			buffer = buffer[decoder.InputOffset():]
		}
	}
}

// 클라이언트로부터 받은 메시지를 Type에 맞게 동작함
func (s *Server) handleMessage(buffer []byte, n int, messages chan messagetype.Message) {
	var message messagetype.Message

	err := json.Unmarshal(buffer[:n], &message)
	if err != nil {
		log.Println("데이터 수신에 실패하였습니다.: ", err.Error())
		return
	}

	log.Println("서버로부터 받은 헤더: ", message.Type)
	messages <- messagetype.Message{
		Type:    message.Type,
		Content: message.Content,
	}
}
