package main

import (
	"fmt"
	"net"
	"os"
)

const bufferSize = 10

type Broker struct {
}

func (b *Broker) HandleClientRequest(requestChan, ackChan chan string) {
	link, _ := net.Listen("tcp", "127.0.0.1:12346")
	defer link.Close()

	for {
		conn, _ := link.Accept()

		go func() {
			buffer := make([]byte, 64)

			_, err := conn.Read(buffer)
			if err == nil && len(requestChan) < bufferSize {
				fmt.Println("Received Message: " + string(buffer))
				requestChan <- string(buffer)

				ackMessage := <-ackChan
				conn.Write([]byte(ackMessage))
			}
		}()
	}
}

func (b *Broker) HandleServerResponse(requestChan, ackChan chan string) {
	for {
		request := <-requestChan

		// Connect
		var err error
		conn, err := net.Dial("tcp", "127.0.0.1:12345")
		if err != nil {
			fmt.Println("failed to connect")
			return
		}
		defer conn.Close()

		// Send Message
		conn.Write([]byte(request))

		go func() {
			buffer := make([]byte, 64)

			_, err = conn.Read(buffer)
			if err == nil {
				fmt.Println("Received Message: " + string(buffer))
				ackChan <- string(buffer)
			}
		}()
	}
}

func (b *Broker) HandleServerRequest(requestChan, ackChan chan string) {
	link, _ := net.Listen("tcp", "127.0.0.1:12347")
	defer link.Close()

	for {
		conn, _ := link.Accept()

		go func() {
			buffer := make([]byte, 64)

			_, err := conn.Read(buffer)
			if err == nil && len(requestChan) < bufferSize {
				fmt.Println("Received Message: " + string(buffer))
				requestChan <- string(buffer)

				ackMessage := <-ackChan
				conn.Write([]byte(ackMessage))
			}
		}()
	}
}

func (b *Broker) HandleClientResponse(requestChan, ackChan chan string) {
	for {
		request := <-requestChan

		// Connect
		var err error
		conn, err := net.Dial("tcp", "127.0.0.1:12344")
		if err != nil {
			fmt.Println("failed to connect")
			return
		}
		defer conn.Close()

		// Send Message
		conn.Write([]byte(request))

		go func() {
			buffer := make([]byte, 64)

			_, err = conn.Read(buffer)
			if err == nil {
				fmt.Println("Received Message: " + string(buffer))
				ackChan <- string(buffer)
			}
		}()
	}
}

func (b *Broker) Run() {
	clientRequestChan := make(chan string, bufferSize)
	serverResponseChan := make(chan string, bufferSize)
	serverRequestChan := make(chan string, bufferSize)
	clientResponseChan := make(chan string, bufferSize)

	go b.HandleClientRequest(clientRequestChan, serverResponseChan)

	go b.HandleServerResponse(clientRequestChan, serverResponseChan)

	go b.HandleServerRequest(serverRequestChan, clientResponseChan)

	go b.HandleClientResponse(serverRequestChan, clientResponseChan)
}

func main() {
	broker := Broker{}
	broker.Run()

	signalChan := make(chan os.Signal)
	forever := make(chan bool)
	go func() {
		<-signalChan
		forever <- true
	}()

	<-forever
}
