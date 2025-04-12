package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Ouacshaman/httpfromtcp/cmd/httpserver/server"
	"github.com/Ouacshaman/httpfromtcp/internal/request"
	"github.com/Ouacshaman/httpfromtcp/internal/response"
)

const port = 42069

func main() {
	server, err := server.Serve(port, proxyHttpbinHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handlerConn(w io.Writer, req *request.Request) {
	writer := response.Writer{
		W:                w,
		StatusCodeWriter: response.StatusWriteSL,
	}

	req.Headers["Connection"] = "close"
	req.Headers["Content-Type"] = "text/html"
	if req.RequestLine.RequestTarget == "/yourproblem" {
		htmlResponse := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
		w.Write([]byte(htmlResponse))
		req.Status = response.BadRq
		req.Headers["Content-Length"] = strconv.Itoa(len(htmlResponse))
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		htmlResponse := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
		req.Status = response.InternalErr
		w.Write([]byte(htmlResponse))
		req.Headers["Content-Length"] = strconv.Itoa(len(htmlResponse))
		return
	}
	okStatus := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
	w.Write([]byte(okStatus))
	req.Status = response.Ok
	req.Headers["Content-Length"] = strconv.Itoa(len(okStatus))
	return

	for writer.StatusCodeWriter != response.StatusComplete {
		switch writer.StatusCodeWriter {
		case response.StatusWriteSL:
			err := writer.WriteStatusLine(req.Status)
			if err != nil {
				fmt.Println(err)
				return
			}
			writer.StatusCodeWriter = response.StatusWriteHeader
		case response.StatusWriteHeader:
			err := writer.WriteHeaders(req.Headers)
			if err != nil {
				fmt.Println(err)
				return
			}
			writer.StatusCodeWriter = response.StatusWriteBody

		case response.StatusWriteBody:
			_, err := writer.WriteBody(b.Bytes())
			if err != nil {
				fmt.Println(err)
				return
			}

			writer.StatusCodeWriter = response.StatusComplete
		default:
			return
		}
	}
	req.Headers["Connection"] = "close"
	req.Headers["Content-Type"] = "text/html"
	if req.RequestLine.RequestTarget == "/yourproblem" {
		htmlResponse := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
		w.Write([]byte(htmlResponse))
		req.Status = response.BadRq
		req.Headers["Content-Length"] = strconv.Itoa(len(htmlResponse))
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		htmlResponse := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
		req.Status = response.InternalErr
		w.Write([]byte(htmlResponse))
		req.Headers["Content-Length"] = strconv.Itoa(len(htmlResponse))
		return
	}
	okStatus := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
	w.Write([]byte(okStatus))
	req.Status = response.Ok
	req.Headers["Content-Length"] = strconv.Itoa(len(okStatus))
	return
}

func defaultResponseHandler(w io.Writer, req *request.Request) {
	req.Headers["Connection"] = "close"
	req.Headers["Content-Type"] = "text/html"
	if req.RequestLine.RequestTarget == "/yourproblem" {
		htmlResponse := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
		w.Write([]byte(htmlResponse))
		req.Status = response.BadRq
		req.Headers["Content-Length"] = strconv.Itoa(len(htmlResponse))
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		htmlResponse := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
		req.Status = response.InternalErr
		w.Write([]byte(htmlResponse))
		req.Headers["Content-Length"] = strconv.Itoa(len(htmlResponse))
		return
	}
	okStatus := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
	w.Write([]byte(okStatus))
	req.Status = response.Ok
	req.Headers["Content-Length"] = strconv.Itoa(len(okStatus))
	return

}

func proxyHttpbinHandler(w io.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	if strings.HasPrefix(target, "/httpbin/") {
		target = strings.TrimPrefix(target, "/httpbin/")
	} else {
		fmt.Println("Does not have /httpbin/ prefix")
		return
	}

	writer := response.Writer{
		W:                w,
		StatusCodeWriter: response.StatusWriteSL,
	}

	url := fmt.Sprintf("https://httpbin.org/%s", target)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	req.Status = response.StatusCode(resp.StatusCode)

	req.Headers = make(map[string]string)
	for k, v := range resp.Header {
		if strings.ToLower(k) != "content-length" {
			req.Headers[k] = v[0]
		}
	}

	req.Headers["Transfer-Encoding"] = "chunked"

	buf := make([]byte, 1028)
	for writer.StatusCodeWriter != response.StatusComplete {
		switch writer.StatusCodeWriter {
		case response.StatusWriteSL:
			err = writer.WriteStatusLine(req.Status)
			if err != nil {
				fmt.Println(err)
				return
			}
			writer.StatusCodeWriter = response.StatusWriteHeader
		case response.StatusWriteHeader:
			err = writer.WriteHeaders(req.Headers)
			if err != nil {
				fmt.Println(err)
				return
			}
			writer.StatusCodeWriter = response.StatusWriteBody
		case response.StatusWriteBody:

			for {
				n, err := resp.Body.Read(buf)
				if n > 0 {
					_, err := writer.WriteChunkedBody(buf[:n])
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				if err != nil {
					if err != io.EOF {
						fmt.Println(err)
					}
					break
				}
			}

			_, err = writer.WriteChunkedBodyDone()
			if err != nil {
				fmt.Println(err)
				return
			}
			writer.StatusCodeWriter = response.StatusComplete
		default:
			return
		}
	}
}
