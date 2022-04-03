package main

import (
	"fmt"
	"net"
	"os"
)

const bufferSize = 10

type Broker struct {
}

func (b *Broker) HandleClientSide(requestChan, ackChan chan string) {
	link, _ := net.Listen("tcp", "127.0.0.1:12346")
	defer link.Close()

	for {
		conn, _ := link.Accept()

		buffer := make([]byte, 64)

		_, err := conn.Read(buffer)
		if err == nil && len(requestChan) < bufferSize {
			fmt.Println("Received Message: " + string(buffer))
			requestChan <- string(buffer)

			go func() {
				ackMessage := <-ackChan
				conn.Write([]byte(ackMessage))
			}()
		}
	}
}

func (b *Broker) HandleServerSide(requestChan, ackChan chan string) {
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

func (b *Broker) Run() {
	requestChan := make(chan string, bufferSize)
	ackChan := make(chan string, bufferSize)

	go b.HandleClientSide(requestChan, ackChan)

	go b.HandleServerSide(requestChan, ackChan)
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
