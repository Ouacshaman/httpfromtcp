package server

import (
	"io"

	"github.com/Ouacshaman/httpfromtcp/internal/request"
)

type HandlerError struct {
	Code    int
	Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
