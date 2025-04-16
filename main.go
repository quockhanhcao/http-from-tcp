package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	linesChannel := getLinesChannel(file)
	for line := range linesChannel {
		fmt.Printf("read: %s\n", line)
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
