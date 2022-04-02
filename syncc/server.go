package main1

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:12345")
	defer conn.Close()

	buffer := make([]byte, 64)
	sequence_number := 0

	for {
		msg := "Message from server, seq_num: " + fmt.Sprint(sequence_number)
		conn.Write([]byte(msg))
		fmt.Println("Sent message: " + msg)

		conn.Read(buffer)
		fmt.Println("Received ack: " + string(buffer))
		conn.Read(buffer)

		time.Sleep(1 * time.Second)
		sequence_number += 1
	}

}
