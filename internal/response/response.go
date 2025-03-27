package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Ouacshaman/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	ok          StatusCode = 200
	badRq       StatusCode = 400
	internalErr StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case ok:
		okStatus := "HTTP/1.1 200 OK"
		_, err := w.Write([]byte(okStatus))
		if err != nil {
			return err
		}
		return nil
	case badRq:
		badRqStatus := "HTTP/1.1 400 Bad Request"
		_, err := w.Write([]byte(badRqStatus))
		if err != nil {
			return err
		}
		return nil
	case internalErr:
		intErrStatus := "HTTP/1.1 500 Internal Server Error"
		_, err := w.Write([]byte(intErrStatus))
		if err != nil {
			return err
		}
		return nil

	default:
		_, err := w.Write([]byte(""))
		if err != nil {
			return err
		}
		return nil
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := make(headers.Headers)
	header["Content-Length"] = strconv.Itoa(contentLen)
	header["Connection"] = "close"
	header["Content-Type"] = "text/plain"
	return header
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	res := ""
	for k, v := range headers {
		header := fmt.Sprintf("%s: %s\n", k, v)
		res += header
	}

	_, err := w.Write([]byte(res))
	if err != nil {
		return err
	}
	return nil
}
