package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	// crlf index
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	// if data starts with crlf
	if idx == 0 {
		return 0, true, nil
	}

	trimedCrlf := data[:idx]

	colonIdx := bytes.Index(trimedCrlf, []byte(":"))
	if colonIdx == -1 {
		return 0, false, fmt.Errorf("Invalid Format: Colon not found")
	}
	if colonIdx > 0 && trimedCrlf[colonIdx-1] == ' ' {
		return 0, false, fmt.Errorf("Invalid Format: Whitespace between Field Name and Colon")
	}

	fieldVal := string(bytes.TrimSpace(trimedCrlf[colonIdx+1:]))
	fieldNameStr := string(bytes.TrimSpace(trimedCrlf[:colonIdx]))
	h[fieldNameStr] = fieldVal

	// +2 for \r\n
	return idx + 2, false, nil
}
