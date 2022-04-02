package main1

import (
	"net"
	"sync"
	"time"
)

type Status int

const (
	READING Status = iota
	PENDING
)

var messageChannel chan []byte
var ackChannel chan string
var wg sync.WaitGroup

func receive() {
	defer wg.Done()
	buffer := make([]byte, 64)
	status := PENDING

	link, _ := net.Listen("tcp", "127.0.0.1:12345")
	defer link.Close()

	for {
		conn, _ := link.Accept()
		defer conn.Close()
		status = READING

		for status == READING {
			_, err := conn.Read(buffer)
			if err != nil {
				status = PENDING
			} else {
				messageChannel <- buffer
				ack := <-ackChannel
				conn.Write([]byte(ack))
			}
		}
	}

}

func send() {
	defer wg.Done()
	for {
		conn, err := net.Dial("tcp", "127.0.0.1:12346")
		for err != nil {
			time.Sleep(1 * time.Second)
			conn, err = net.Dial("tcp", "127.0.0.1:12346")
		}
		defer conn.Close()
		buffer := <-messageChannel
		conn.Write(buffer)
		ackChannel <- "Ack for: " + string(buffer)
		conn.Read(buffer)
	}
}

func main() {
	messageChannel = make(chan []byte)
	ackChannel = make(chan string)
	wg.Add(2)

	go receive()
	go send()
	wg.Wait()
}
