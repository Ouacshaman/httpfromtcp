package request

import (
	"errors"
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

func RequestFromReader(reader io.Reader) (*Request, error) {
	var request Request
	b := make([]byte, 64)
	_, err := reader.Read(b)
	if err != nil {
		return &request, err
	}

	lines := strings.Split(string(b), "\r\n")

	rq_line := lines[0]
	rq_parts := strings.Split(rq_line, " ")

	if len(rq_parts) != 3 {
		fmt.Println("Invalid Amounts of Parts on the Request Line")
		return &request, errors.New("Invalid Amount of Parts on Request Line")
	}

	if rq_parts[0] != "POST" && rq_parts[0] != "GET" {
		fmt.Println("Invalid Placement and not a valid method")
		return &request, errors.New("Invalid Method and Placement")
	}

	if rq_parts[0] != strings.ToUpper(rq_parts[0]) {
		fmt.Println("Method not capitalized: ", rq_parts[0])
		return &request, errors.New("Method not capitalized")
	}

	if strings.HasPrefix(rq_parts[1], "/") == false {
		fmt.Println("Invalid Request Target")
		return &request, errors.New("Invalid Request Target")
	}

	if rq_parts[2] != "HTTP/1.1" {
		fmt.Printf("The version is incorrect and not HTTP/1.1: %s", rq_parts[2])
		return &request, errors.New("Incorrect HTTP version")
	}

	request.RequestLine.Method = rq_parts[0]
	request.RequestLine.RequestTarget = rq_parts[1]
	request.RequestLine.HttpVersion = strings.TrimPrefix(rq_parts[2], "HTTP/")

	return &request, nil
}
