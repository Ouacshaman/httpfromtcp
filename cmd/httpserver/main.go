package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ouacshaman/httpfromtcp/cmd/httpserver/server"
	"github.com/Ouacshaman/httpfromtcp/internal/request"
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

func handlerConn(w io.Writer, req *request.Request) *server.HandlerError {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handlerErr := server.HandlerError{
			Code:    400,
			Message: "Your problem is not my problem\n",
		}
		return &handlerErr
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handlerErr := server.HandlerError{
			Code:    500,
			Message: "Woopsie, my bad\n",
		}
		return &handlerErr
	}
	okStatus := "All good, frfr\n"
	w.Write([]byte(okStatus))
	return nil
}
