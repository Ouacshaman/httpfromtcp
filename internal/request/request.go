package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/go-faster/errors"
)

type State int

const (
	Initialized State = iota
	Done
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	State
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	rqLine, err := parseRequestLine(raw)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rqLine,
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.RequestLine.State {
	case Initialized:
		numByte, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if numByte == 0 {
			return 0, nil
		}
		r.RequestLine.State = Done
		return numByte, nil
	case Done:
		return 0, fmt.Errorf("error: trying to read in Done State: %d", Done)
	default:
		return 0, errors.New("Unknow State")
	}
}

func parseRequestLine(data []byte) (int, err) {
	crlfInd := bytes.Index(data, []byte(crlf))
	if crlfInd == -1 {
		return 0
	}
	_, err := requestLineFromString(string(data[:crlfInd]))
	if err != nil {
		return nil, err
	}
	return len(data[:crlfInd]), nil
}

func requestLineFromString(rqLine string) (*RequestLine, error) {
	rqParts := strings.Split(rqLine, " ")

	if len(rqParts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", rqLine)
	}

	for _, v := range rqParts[0] {
		if v < 'A' || v > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", rqParts[0])
		}
	}

	versionParts := strings.Split(rqParts[2], "/")

	if versionParts[0] != "HTTP" {
		return nil, fmt.Errorf("invalid Http version: %s", versionParts[0])
	}

	if versionParts[1] != "1.1" {
		return nil, fmt.Errorf("invalid Http version: %s", versionParts[0])
	}

	return &RequestLine{
		Method:        rqParts[0],
		RequestTarget: rqParts[1],
		HttpVersion:   versionParts[1],
	}, nil

}
