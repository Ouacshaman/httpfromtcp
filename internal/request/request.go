package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
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

func parseRequestLine(data []byte) (*RequestLine, error) {
	crlfInd := bytes.Index(data, []byte(crlf))
	if crlfInd == -1 {
		return nil, fmt.Errorf("CRLF not found in request-line")
	}
	requestLine, err := requestLineFromString(string(data[:crlfInd]))
	if err != nil {
		return nil, err
	}
	return requestLine, nil
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
