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
	w := response.Writer{
		W:                conn,
		StatusCodeWriter: response.StatusWriteSL,
	}
	rq, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteError(response.BadRq, err.Error())
		return
	}
	var b bytes.Buffer
	s.handler(&b, rq)
	//header := response.GetDefaultHeaders(len(b.Bytes()))
	for w.StatusCodeWriter != response.StatusComplete {
		switch w.StatusCodeWriter {
		case response.StatusWriteSL:
			err := w.WriteStatusLine(response.Ok)
			if err != nil {
				fmt.Println(err)
			}
			w.StatusCodeWriter = response.StatusWriteHeader
		case response.StatusWriteHeader:
			header := response.GetDefaultHeaders(len(b.Bytes()))
			err := w.WriteHeaders(header)
			if err != nil {
				fmt.Println(err)
			}
			w.StatusCodeWriter = response.StatusWriteBody

		case response.StatusWriteBody:
			w.WriteBody(b.Bytes())
			w.StatusCodeWriter = response.StatusComplete
		default:
			return
		}
	}
}
