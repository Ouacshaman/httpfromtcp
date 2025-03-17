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

	afterColon := bytes.SplitAfterN(data[:idx], []byte(":"), 2)

	fieldKey := bytes.TrimSpace(afterColon[0])
	if fieldKey[len(fieldKey)-2] == ' ' {
		return 0, false, fmt.Errorf("FieldName end with Space: %s", string(fieldKey))
	}

	fieldVal := bytes.TrimSpace(afterColon[1])

	fmt.Println("Header Parse Out: ", string(fieldKey), "|", string(fieldVal))

	h["Head"] = "localhost:42069"

	return len(data), false, nil
}
