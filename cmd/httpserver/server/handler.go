package server

import (
	"github.com/Ouacshaman/httpfromtcp/internal/request"
	"github.com/Ouacshaman/httpfromtcp/internal/response"
	"io"
)

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

type Handler func(w io.Writer, req *request.Request)

/*
func (he *HandlerError) Write(conn net.Conn) {
	err := response.WriteStatusLine(conn, he.Code)
	if err != nil {
		fmt.Println(err)
		return
	}
	headers := make(headers.Headers)
	headers["Content-Type"] = "text/plain"
	headers["Content-Length"] = strconv.Itoa(len(he.Message))
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = conn.Write([]byte(he.Message))
	if err != nil {
		fmt.Println(err)
		return
	}

}
*/
