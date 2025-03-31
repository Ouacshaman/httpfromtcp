package server

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/Ouacshaman/httpfromtcp/internal/headers"
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
	rq, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("Unable to read request: %v", err)
		return
	}
	var b bytes.Buffer
	handlerErr := s.handler(&b, rq)
	if handlerErr != nil {
		err := response.WriteStatusLine(conn, response.StatusCode(handlerErr.Code))
		if err != nil {
			fmt.Println(err)
			return
		}
		headers := make(headers.Headers)
		headers["Content-Type"] = "text/plain"
		headers["Content-Length"] = strconv.Itoa(len(handlerErr.Message))
		err = response.WriteHeaders(conn, headers)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = conn.Write([]byte(handlerErr.Message))
		if err != nil {
			fmt.Println(err)
			return
		}
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

	defer conn.Close()
}
