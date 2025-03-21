package request

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Ouacshaman/httpfromtcp/internal/headers"
	"io"
	"strings"
)

type State int

const (
	requestStateInitialized State = iota
	requestStateDone
	requestStateParsingHeaders
)

type Request struct {
	State
	RequestLine RequestLine
	Headers     headers.Headers
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	bufferSize := 8
	buf := make([]byte, bufferSize)
	readToIndex := 0
	rq := &Request{State: requestStateInitialized}
	rq.Headers = make(headers.Headers)

	for rq.State != requestStateDone {
		if readToIndex >= len(buf) {
			dbl := make([]byte, len(buf)*2)
			copy(dbl, buf[:readToIndex])
			buf = dbl
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil && errors.Is(err, io.EOF) {
			return nil, err
		}
		readToIndex += n

		bytesPassed, err := rq.parse(buf[:readToIndex])
		if err != nil {
			return nil, errors.New("Unable to Parse into buffer")
		}

		/*
			New slice created and moved to the beginning of buffer,
			which removed data that was already parsed.
			The decrement allow our data to start at the end of the moved data.
		*/
		if bytesPassed > 0 {
			copy(buf, buf[bytesPassed:readToIndex])
			readToIndex -= bytesPassed
		}

		if err == io.EOF {
			if rq.State != requestStateDone && readToIndex > 0 {
				return nil, errors.New("Unexepected EOF before fulling parsing data")
			}
			break
		}

	}
	return rq, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	i := 0
	max := 25
	for r.State != requestStateDone {
		i++
		if i > max {
			break
		}
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return totalBytesParsed, err
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case requestStateInitialized:
		numByte, err := r.parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if numByte == 0 {
			return 0, nil
		}
		r.State = requestStateParsingHeaders
		return numByte, nil
	case requestStateParsingHeaders:
		bytesParsed, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = requestStateDone
		}
		return bytesParsed, nil
	default:
		return 0, errors.New("Unknow State")
	}
}
func (r *Request) parseRequestLine(data []byte) (int, error) {
	crlfInd := bytes.Index(data, []byte(crlf))
	if crlfInd == -1 {
		return 0, nil
	}
	rqLineStr := string(data[:crlfInd])
	err := r.requestLineFromString(rqLineStr)
	if err != nil {
		return 0, err
	}
	return crlfInd + len(crlf), nil
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
	r.RequestLine.HttpVersion = versionParts[1]

	return nil

}
