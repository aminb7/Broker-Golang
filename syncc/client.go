package main1

import (
	"fmt"
	"net"
)

func main() {
	link, _ := net.Listen("tcp", "127.0.0.1:12346")
	defer link.Close()

	buffer := make([]byte, 64)

	for i := 0; i < 5; i++ {
		conn, _ := link.Accept()
		defer conn.Close()
		_, err := conn.Read(buffer)
		if err == nil {
			fmt.Println(string(buffer))
			conn.Write([]byte("OK"))
		}
	}

}
