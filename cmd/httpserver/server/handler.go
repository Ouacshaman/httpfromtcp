package server

import (
	"io"

	"github.com/Ouacshaman/httpfromtcp/internal/request"
)

type Handler func(w io.Writer, req *request.Request)
