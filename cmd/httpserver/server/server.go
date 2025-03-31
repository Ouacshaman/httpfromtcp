package server

import (
	"bytes"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/Ouacshaman/httpfromtcp/internal/request"
	"github.com/Ouacshaman/httpfromtcp/internal/response"
)

type Server struct {
	ln      net.Listener
	closed  *atomic.Bool
	handler Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	stringPort := fmt.Sprintf(":%d", port)
	listen, err := net.Listen("tcp", stringPort)
	if err != nil {
		return nil, err
	}

	server := &Server{
		ln:      listen,
		closed:  &atomic.Bool{},
		handler: handler,
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
	defer conn.Close()
	rq, err := request.RequestFromReader(conn)
	if err != nil {
		handleErr := &HandlerError{
			Code:    response.Ok,
			Message: err.Error(),
		}
		handleErr.Write(conn)
		return
	}
	var b bytes.Buffer
	handlerErr := s.handler(&b, rq)
	if handlerErr != nil {
		handlerErr.Write(conn)
		return
	}
	header := response.GetDefaultHeaders(len(b.Bytes()))
	err = response.WriteStatusLine(conn, 200)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = response.WriteHeaders(conn, header)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn.Write(b.Bytes())
}
