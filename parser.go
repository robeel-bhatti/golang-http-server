package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
}

func ParseRequest(r *bufio.Reader) (*Request, error) {
	// keep reading from bufio's internal buffer until it finds the delimiter.
	// it will return everything up to and including the delimiter.
	// if the reader does not find the delimiter in the bufio internal buffer,
	// it will pull from the socket buffer. If the delimiter is still not found, it will block
	// the current goroutine until the delimiter is found.
	rl, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	rl = strings.TrimRight(rl, "\r\n")
	parts := strings.Split(rl, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("malformed request line: %s", rl)
	}
	req := &Request{Method: parts[0], Path: parts[1], Version: parts[2]}

	h := make(map[string]string)
	for {
		hl, err := r.ReadString('\n')
		if err != nil {
			// caller should check for IOF here in case
			// this is the last h of the HTTP request
			// and no http body was provided
			return nil, err
		}

		hl = strings.TrimRight(hl, "\r\n")
		if hl == "" {
			break
		}

		idx := strings.Index(hl, ":")
		if idx == -1 {
			return nil, fmt.Errorf("invalid h: %v", hl)
		}
		one := strings.ToLower(strings.TrimSpace(hl[:idx]))
		two := strings.TrimSpace(hl[idx+1:])
		h[one] = two
	}

	req.Headers = h
	cl, ok := req.Headers["content-length"]
	if !ok {
		return req, nil
	}

	n, err := strconv.Atoi(cl)
	if err != nil || n < 0 {
		return nil, fmt.Errorf("invalid content-length: %v", cl)
	}

	body := make([]byte, n)
	_, err = io.ReadFull(r, body)
	if err != nil {
		return nil, err
	}

	req.Body = body
	return req, nil
}
