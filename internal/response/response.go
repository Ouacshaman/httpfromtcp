package response

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/Ouacshaman/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Ok          StatusCode = 200
	BadRq       StatusCode = 400
	InternalErr StatusCode = 500
)

type WriterStatusCode int

const (
	StatusWriteSL WriterStatusCode = iota
	StatusWriteHeader
	StatusWriteBody
	StatusComplete
)

type Writer struct {
	W                io.Writer
	StatusCodeWriter WriterStatusCode
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	switch statusCode {
	case Ok:
		okStatus := "HTTP/1.1 200 OK\r\n"
		_, err := w.W.Write([]byte(okStatus))
		if err != nil {
			return err
		}
		return nil
	case BadRq:
		badRqStatus := "HTTP/1.1 400 Bad Request\r\n"
		_, err := w.W.Write([]byte(badRqStatus))
		if err != nil {
			return err
		}
		return nil
	case InternalErr:
		intErrStatus := "HTTP/1.1 500 Internal Server Error\r\n"
		_, err := w.W.Write([]byte(intErrStatus))
		if err != nil {
			return err
		}
		return nil

	default:
		_, err := w.W.Write([]byte(""))
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

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	res := ""
	for k, v := range headers {
		header := fmt.Sprintf("%s: %s\r\n", k, v)
		res += header
	}

	res += "\r\n"

	_, err := w.W.Write([]byte(res))
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.W.Write(p)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) WriteError(code StatusCode, message string) {
	err := w.WriteStatusLine(code)
	if err != nil {
		fmt.Println(err)
		return
	}
	headers := make(headers.Headers)
	headers["Content-Type"] = "text/plain"
	headers["Content-Length"] = strconv.Itoa(len(message))
	err = w.WriteHeaders(headers)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = w.W.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	temp := p
	idx := bytes.Index(temp, []byte("\r\n"))
	hexN, err := strconv.ParseInt(string(p[:idx]), 16, 64)
	if err != nil {
		return 0, err
	}

	w.W.Write(temp[idx : idx+int(hexN)])

	return int(hexN) + 2 + len(temp[:idx]), nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.W.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return 0, err
	}
	return n, nil
}
