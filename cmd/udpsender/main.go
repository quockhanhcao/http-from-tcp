package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// an udp address is a struct with IP, Port, and Zone
	udpAddress, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving UDP address: %v\n", err)
		os.Exit(1)
	}
	udpConnection, err := net.DialUDP("udp", nil, udpAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dialing UDP: %v\n", err)
		os.Exit(1)
	}
	defer udpConnection.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		switch input, err := reader.ReadString('\n'); err {
		// if the read succeeded, write the input to udp connection
		case nil:
			if _, err = udpConnection.Write([]byte(input)); err != nil {
				fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Message sent: %s", input)
		// the end of the input, exit gracefully
		case io.EOF:
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
	}
}
