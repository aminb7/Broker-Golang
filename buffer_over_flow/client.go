package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

type Client struct {
}

func (c *Client) readInputMessage() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Message: ")
	message, _ := reader.ReadString('\n')

	return message
}

func (c *Client) Run() {
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
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	client := Client{}
	client.Run()
}
