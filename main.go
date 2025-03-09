package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("Unable to Listen to TCP port: %v\n", err)
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Unable to accept connection: %v\n", err)
			return
		}
		fmt.Println("Connection have been accepted")
		go func() {
			defer conn.Close()
			ch := getLinesChannel(conn)
			/*
				The loop below is similar to
				for {
				    v, ok := <-ch
				    if !ok {
				        break
				    }
				    fmt.Printf("%s\n", v)
				}
				But it is more concise, less error-prone, and iterates the values similary in a more clear way
				this will iterate till the channel is closed
			*/
			for line := range ch {
				fmt.Printf("%s\n", line)
			}
			fmt.Println("Connection Closed")
		}()
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		/*
			channel is closed inside as the go routine is the sender
			The goroutine handles everything related to sending values
			to the channel.
			It doesn't rely on external code to manage the channel's lifecycle.
			This will ensure channel does not close prematurely
			This will ensure the channel close after the go routine is processed and no data is sent to a closed channel
		*/
		defer close(ch)

		res := ""

		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err != nil {
				break
			}

			parts := strings.Split(string(b[:n]), "\n")

			for _, v := range parts[:len(parts)-1] {
				res += v
				ch <- res
				res = ""
			}
			res += parts[len(parts)-1]
		}
		ch <- res
	}()
	return ch
}
