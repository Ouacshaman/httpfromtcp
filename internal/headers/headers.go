package headers

import (
	"bytes"
	"fmt"
	"strings"
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
		return 2, true, nil
	}

	trimedCrlf := data[:idx]

	colonIdx := bytes.Index(trimedCrlf, []byte(":"))
	if colonIdx == -1 {
		return 0, false, fmt.Errorf("Invalid Format: Colon not found")
	}
	if colonIdx > 0 && trimedCrlf[colonIdx-1] == ' ' {
		return 0, false, fmt.Errorf("Invalid Format: Whitespace between Field Name and Colon")
	}

	special := []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}
	fieldVal := string(bytes.TrimSpace(trimedCrlf[colonIdx+1:]))
	fieldNameStr := strings.ToLower(string(bytes.TrimSpace(trimedCrlf[:colonIdx])))
	for _, v := range []byte(fieldNameStr) {
		letters := false
		numbers := false
		sc := false
		if 'a' <= v && 'z' >= v {
			letters = true
		}
		if '0' <= v && '9' >= v {
			numbers = true
		}
		for _, n := range special {
			if v == n {
				sc = true
			}
		}
		if (numbers || letters || sc) == false {
			return 0, false, fmt.Errorf("Field Name does not match requirements")
		}
	}
	elem, ok := h[fieldNameStr]
	if !ok {
		h[fieldNameStr] = fieldVal
		return idx + 2, false, nil
	}

	h[fieldNameStr] = strings.Join([]string{elem, fieldVal}, ", ")
	// +2 for \r\n
	return idx + 2, false, nil
}

func (h Headers) Get(key string) string {
	lowered := strings.ToLower(key)
	return h[lowered]
}
