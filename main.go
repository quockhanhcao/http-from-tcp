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
	defer file.Close()

	currentLine := ""
	for {
		buffer := make([]byte, 8, 8)
		n, err := file.Read(buffer)
		if err != nil {
			if currentLine != "" {
				fmt.Printf("read: %s\n", currentLine)
				currentLine = ""
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
            fmt.Printf("read: %s\n", currentLine)
			currentLine = ""
		}
		currentLine += parts[len(parts) - 1]
	}
}
