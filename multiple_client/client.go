package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
}

func (c *Client) readInputMessage() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Message: ")
	message, _ := reader.ReadString('\n')

	return message
}

func (c *Client) HandleRequest() {
	for {
		// Connect
		var err error
		conn, err := net.Dial("tcp", "127.0.0.1:12346")
		if err != nil {
			fmt.Println("failed to connect")
			return
		}
		defer conn.Close()

		// Send Message
		message := c.readInputMessage()
		conn.Write([]byte(message))

		buffer := make([]byte, 64)

		go func() {
			_, err = conn.Read(buffer)
			if err == nil {
				fmt.Println("Received Message: " + string(buffer))
			}
		}()

	}
}

func (c *Client) HandleResponse() {
	link, _ := net.Listen("tcp", "127.0.0.1:12344")
	defer link.Close()

	for {
		conn, _ := link.Accept()
		defer conn.Close()

		buffer := make([]byte, 64)

		_, err := conn.Read(buffer)
		if err == nil {
			fmt.Println("Received Message: " + string(buffer))
			conn.Write([]byte("Client Ack: " + string(buffer)))
		}
	}
}

func (c *Client) Run() {
	go c.HandleRequest()
	go c.HandleResponse()
}

func main() {
	client := Client{}
	client.Run()

	signalChan := make(chan os.Signal)
	forever := make(chan bool)
	go func() {
		<-signalChan
		forever <- true
	}()

	<-forever
}
