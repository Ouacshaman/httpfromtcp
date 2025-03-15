package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type State int

const (
	Initialized State = iota
	Done
)

type Request struct {
	State
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader, bufferSize int) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)

	readToIndex := 0

	var rq Request

	rq.State = Initialized

	for rq.State != Done {
		if readToIndex == len(buf)-1 {
			dbl := make([]byte, len(buf)*2, len(buf)*2)
			copy(dbl, buf)
			buf = dbl
		}
		for {
			n, err := reader.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					rq.State = Done
					break
				}

				fmt.Println("Error Found: ", err)
				return nil, err
			}
			fmt.Println(string(buf), n)
			readToIndex = n
			n, err = rq.parse(buf[:n])
			if err != nil {
				return nil, errors.New("Unable to Parse into buffer")
			}
			empty := make([]byte, 0)
			copy(buf, empty)
			readToIndex -= n
		}
	}
	return &rq, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		numByte, err := r.parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if numByte == 0 {
			return 0, nil
		}
		r.State = Done
		return numByte, nil
	case Done:
		return 0, fmt.Errorf("error: trying to read in Done State: %d", Done)
	default:
		return 0, errors.New("Unknow State")
	}
}

func (r *Request) parseRequestLine(data []byte) (int, error) {
	crlfInd := bytes.Index(data, []byte(crlf))
	if crlfInd == -1 {
		return 0, nil
	}
	err := r.requestLineFromString(string(data[:crlfInd]))
	if err != nil {
		return 0, err
	}
	return len(data[:crlfInd]), nil
}

func (r *Request) requestLineFromString(rqLine string) error {
	rqParts := strings.Split(rqLine, " ")

	if len(rqParts) != 3 {
		return fmt.Errorf("poorly formatted request-line: %s", rqLine)
	}

	for _, v := range rqParts[0] {
		if v < 'A' || v > 'Z' {
			return fmt.Errorf("invalid method: %s", rqParts[0])
		}
	}

	versionParts := strings.Split(rqParts[2], "/")

	if versionParts[0] != "HTTP" {
		return fmt.Errorf("invalid Http version: %s", versionParts[0])
	}

	if versionParts[1] != "1.1" {
		return fmt.Errorf("invalid Http version: %s", versionParts[0])
	}

	r.RequestLine.Method = rqParts[0]
	r.RequestLine.RequestTarget = rqParts[1]
	r.RequestLine.HttpVersion = rqParts[2]

	return nil

}
