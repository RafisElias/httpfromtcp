package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		for line := range getLisnesChannel(conn) {
			fmt.Println(line)
		}
	}
}

func getLisnesChannel(f io.ReadCloser) <-chan string {
	channel := make(chan string)

	go func() {
		defer f.Close()
		defer close(channel)
		var currentLine string = ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			idx := bytes.IndexByte(data, '\n')

			if idx != -1 {
				currentLine += string(data[:idx])
				channel <- currentLine
				currentLine = string(data[idx+1 : n])
			} else {
				currentLine += string(data[:n])
			}
		}
		// Send any remaining content when connection closes
		if currentLine != "" {
			channel <- currentLine
		}
	}()

	return channel
}
