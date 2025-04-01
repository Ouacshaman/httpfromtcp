package server

import (
	"io"
	"net"

	"github.com/Ouacshaman/httpfromtcp/internal/request"
	"github.com/Ouacshaman/httpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request)
