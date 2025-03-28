package server

import (
	"io"

	"github.com/Ouacshaman/httpfromtcp/internal/request"
)

type HandlerError struct {
	code    int
	message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
