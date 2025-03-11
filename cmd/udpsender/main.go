package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		fmt.Println("UDP address connection failed")
		return
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("UDP connection attempt failed")
		return
	}
	defer udpConn.Close()

	bufReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")

		line, err := bufReader.ReadString('\n')
		if err != nil {
			fmt.Println("Unable to read line to string")
		}

		udpConn.Write([]byte(line))
	}
}
