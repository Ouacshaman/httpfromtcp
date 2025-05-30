package response

import (
	"fmt"
	"github.com/Ouacshaman/httpfromtcp/internal/headers"
)

func (w *Writer) WriteTrailers(h headers.Headers) error {
	res := ""
	for k, v := range h {
		header := fmt.Sprintf("%s: %s\n", k, v)
		res += header
	}

	res += "\r\n"

	_, err := w.W.Write([]byte(res))
	if err != nil {
		return err
	}
	return nil
}
