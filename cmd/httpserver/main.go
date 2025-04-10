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
	buf := make([]byte, 1028)
	target := req.RequestLine.RequestTarget
	if strings.HasPrefix(target, "/httpbin/") {
		target = strings.TrimPrefix(target, "/httpbin/")
	} else {
		fmt.Println("Does not have /httpbin/ prefix")
		return
	}

	_, ok := req.Headers["Content-Length"]
	if ok {
		delete(req.Headers, "Content-Length")
	}

	req.Headers["Transfer-Encoding"] = "chunked"

	url := fmt.Sprintf("https://httpbin.org/%s", target)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	n, err := resp.Body.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Data read from /httpbin: %d\n", n)
	_, err = w.Write(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
}
