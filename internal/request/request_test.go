package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}

func TestRequestLineParse(t *testing.T) {
	// Test: Good GET Request line
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	reader = &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	reader = &chunkReader{
		data:            "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)
}

func TestParseRequestLine_Complete(t *testing.T) {
	method, target, version, consumed, err := parseRequestLine([]byte("GET / HTTP/1.1\r\n"))
	require.NoError(t, err)
	assert.Equal(t, "GET", method)
	assert.Equal(t, "/", target)
	assert.Equal(t, "1.1", version)
	assert.Equal(t, 16, consumed)
}

func TestParseRequestLine_Incomplete(t *testing.T) {
	_, _, _, consumed, err := parseRequestLine([]byte("GET / HTTP"))
	require.NoError(t, err)
	assert.Equal(t, 0, consumed) // Should need more data
}

func TestRequestParse_StateTransition(t *testing.T) {
	req := &Request{state: stateInitialized}
	consumed, err := req.parse([]byte("GET / HTTP/1.1\r\n"))
	require.NoError(t, err)
	assert.Equal(t, 16, consumed)
	assert.Equal(t, stateDone, req.state)
}

func TestRequestFromReader_Chunked(t *testing.T) {
	reader := &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost\r\n\r\n",
		numBytesPerRead: 3, // Force chunked reading
	}
	req, err := RequestFromReader(reader)
	require.NoError(t, err)
	assert.Equal(t, "GET", req.RequestLine.Method)
	assert.Equal(t, "/coffee", req.RequestLine.RequestTarget)
}
