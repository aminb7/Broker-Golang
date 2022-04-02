package main

import (
	"fmt"
	"net"
)

type Server struct {
}

func (s *Server) Run() {
	link, _ := net.Listen("tcp", "127.0.0.1:12345")
	defer link.Close()

	for {
		conn, _ := link.Accept()
		defer conn.Close()

		buffer := make([]byte, 64)

		_, err := conn.Read(buffer)
		if err == nil {
			fmt.Println("Received Message: " + string(buffer))
			conn.Write([]byte("Ack: " + string(buffer)))
		}
	}
}

func main() {
	server := Server{}
	server.Run()
}
