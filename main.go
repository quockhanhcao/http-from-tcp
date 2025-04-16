package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
		linesChannel := getLinesChannel(connection)
		for line := range linesChannel {
			fmt.Println(line)
		}
		fmt.Println("Connection closed from:", connection.RemoteAddr())
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	channel := make(chan string)
	go func() {
		defer f.Close()
		defer close(channel)
		currentLine := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					channel <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				currentLine += parts[i]
				channel <- currentLine
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()
	return channel
}
