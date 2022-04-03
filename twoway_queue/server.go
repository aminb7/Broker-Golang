package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

type Server struct {
}

func (s *Server) readInputMessage() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Message: ")
	message, _ := reader.ReadString('\n')

	return message
}

func (s *Server) HandleRequest() {
	for {
		// Connect
		var err error
		conn, err := net.Dial("tcp", "127.0.0.1:12347")
		if err != nil {
			fmt.Println("failed to connect")
			return
		}
		defer conn.Close()

		// Send Message
		message := s.readInputMessage()
		conn.Write([]byte(message))

		buffer := make([]byte, 64)

		go func() {
			_, err = conn.Read(buffer)
			if err == nil {
				fmt.Println("Received Message: " + string(buffer))
			}
		}()
		time.Sleep(200 * time.Millisecond)
	}
}

func (s *Server) HandleResponse() {
	link, _ := net.Listen("tcp", "127.0.0.1:12345")
	defer link.Close()

	for {
		conn, _ := link.Accept()
		defer conn.Close()

		buffer := make([]byte, 64)

		_, err := conn.Read(buffer)
		if err == nil {
			fmt.Println("Received Message: " + string(buffer))
			conn.Write([]byte("Server Ack: " + string(buffer)))
		}
	}
}

func (s *Server) Run() {
	go s.HandleRequest()
	go s.HandleResponse()
}

func main() {
	server := Server{}
	server.Run()

	signalChan := make(chan os.Signal)
	forever := make(chan bool)
	go func() {
		<-signalChan
		forever <- true
	}()

	<-forever
}
