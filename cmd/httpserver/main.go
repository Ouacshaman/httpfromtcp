package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Ouacshaman/httpfromtcp/cmd/httpserver/server"
	"github.com/Ouacshaman/httpfromtcp/internal/request"
	"github.com/Ouacshaman/httpfromtcp/internal/response"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handlerConn)
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
