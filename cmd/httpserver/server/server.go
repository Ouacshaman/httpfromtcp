package server

import (
	"fmt"
	"net"
)

type Server struct {
	ln net.Listener
}

func Serve(port int) (*Server, error) {
	stringPort := fmt.Sprintf(":%d", port)
	listen, err := net.Listen("tcp", stringPort)
	if err != nil {
		return nil, err
	}
	return &Server{
		ln: listen,
	}, nil
}
