package headers

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewHeaders() Headers {
	return make(Headers)
}

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header

	headers := NewHeaders()
	data := []byte("Host!#$%&'*+-.^_`|~: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host!#$%&'*+-.^_`|~"])
	assert.Equal(t, 38, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid Spacing Header
	headers = NewHeaders()
	data = []byte("    Hos22t:  localhost:42069  \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["hos22t"])
	assert.Equal(t, 32, n)
	assert.False(t, done)

	// Test: Invalid Header Name
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

}

func TestMultiHeaders(t *testing.T) {
	headers := NewHeaders()
	data := []byte("  Host:   localhost:42069 \r\n  User-Agent:   curl/7.81.0\r\n\r\n")
	list := bytes.SplitAfterN(data, []byte("\r\n"), 3)

	n, done, err := headers.Parse(list[0])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 28, n)
	assert.False(t, done)
	n, done, err = headers.Parse(list[1])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.Equal(t, 29, n)
	assert.False(t, done)
	_, done, _ = headers.Parse(list[2])
	assert.True(t, done)
}
