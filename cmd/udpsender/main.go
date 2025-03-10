package main

import (
	"fmt"
	"net"
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
}
