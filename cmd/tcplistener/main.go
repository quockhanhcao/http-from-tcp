package main

import (
	"fmt"
	"github.com/quockhanhcao/http-from-tcp/internal/request"
	"log"
	"net"
)

func main() {
	tcpListener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Error listening for TCP: %s\n", err.Error())
	}
	defer tcpListener.Close()
	for {
		connection, err := tcpListener.Accept()
		if err != nil {
			log.Fatalf("Error accepting TCP connection: %s\n", err.Error())
		}
		fmt.Println("Connection accepted from:", connection.RemoteAddr())
		request, _ := request.RequestFromReader(connection)
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range request.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Println("Connection closed from:", connection.RemoteAddr())
	}
}
