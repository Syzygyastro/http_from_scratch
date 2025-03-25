package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type parserState int

const (
	stateInitialized parserState = iota
	stateDone
)

type Request struct {
	RequestLine RequestLine
	state       parserState
	buffer      []byte
}

type RequestLine struct {
	HttpVersion   string // Should store just "1.1"
	RequestTarget string
	Method        string
}

func parseRequestLine(data []byte) (method, target, version string, consumed int, err error) {
	end := bytes.Index(data, []byte("\r\n"))
	if end == -1 {
		return "", "", "", 0, nil
	}

	line := string(data[:end])
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return "", "", "", 0, errors.New("invalid request line format")
	}

	// Validate method
	for _, c := range parts[0] {
		if c < 'A' || c > 'Z' {
			return "", "", "", 0, errors.New("invalid HTTP method")
		}
	}

	// Extract version number
	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 || versionParts[0] != "HTTP" {
		return "", "", "", 0, errors.New("invalid HTTP version format")
	}

	return parts[0], parts[1], versionParts[1], end + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == stateDone {
		return 0, nil
	}

	method, target, version, consumed, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}
	if consumed == 0 {
		return 0, nil
	}

	r.RequestLine = RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version, // Now just "1.1"
	}
	r.state = stateDone

	return consumed, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		state:  stateInitialized,
		buffer: make([]byte, 0, 8),
	}

	readBuf := make([]byte, 8)
	for req.state != stateDone {
		n, err := reader.Read(readBuf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if n == 0 {
			break
		}

		req.buffer = append(req.buffer, readBuf[:n]...)
		consumed, err := req.parse(req.buffer)
		if err != nil {
			return nil, err
		}

		if consumed > 0 {
			req.buffer = req.buffer[consumed:]
		}
	}

	if req.state != stateDone {
		return nil, errors.New("incomplete request")
	}

	return req, nil
}
