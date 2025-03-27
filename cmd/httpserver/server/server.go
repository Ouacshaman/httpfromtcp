package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/Ouacshaman/httpfromtcp/internal/response"
)

type Server struct {
	ln     net.Listener
	closed *atomic.Bool
}

func Serve(port int) (*Server, error) {
	stringPort := fmt.Sprintf(":%d", port)
	listen, err := net.Listen("tcp", stringPort)
	if err != nil {
		return nil, err
	}

	server := &Server{
		ln:     listen,
		closed: &atomic.Bool{},
	}

	server.closed.Store(false)
	server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.ln.Close()
}

func (s *Server) listen() {

	closed := &atomic.Bool{}
	closed.Store(false)

	go func() {
		for !closed.Load() {
			conn, err := s.ln.Accept()
			if err != nil {
				if !closed.Load() {
					fmt.Printf("Unable to accept connection: %v\n", err)
				}
				continue
			}

			go s.handle(conn)
		}
	}()
}

func (s *Server) handle(conn net.Conn) {
	_ = response.WriteStatusLine(conn, 200)
	header := response.GetDefaultHeaders(0)
	_ = response.WriteHeaders(conn, header)
	defer conn.Close()
}
