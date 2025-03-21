package main

import (
	"fmt"
	"net"

	"github.com/Ouacshaman/httpfromtcp/internal/request"
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
		fmt.Println("Connection: ", conn.RemoteAddr(), " have been accepted")
		rq, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Unable to generate request from Connection")
			return
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", rq.RequestLine.Method, rq.RequestLine.RequestTarget, rq.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		headers := rq.Headers
		for k, v := range headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
	}
}
