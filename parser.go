package main

import (
	"bufio"
	"fmt"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers []Header
}

type Header struct {
	Key   string
	Value string
}

func ParseRequest(r *bufio.Reader) (*Request, error) {
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

	for {
		hl, err := r.ReadString('\n')
		if err != nil {
			// caller should check for IOF here in case
			// this is the last header of the HTTP request
			// and no http body was provided
			return nil, err
		}

		hl = strings.TrimRight(hl, "\r\n")
		if hl == "" {
			break
		}

		idx := strings.Index(hl, ":")
		if idx == -1 {
			return nil, fmt.Errorf("invalid header: %v", hl)
		}
		partOne := strings.TrimSpace(hl[:idx])
		partTwo := strings.TrimSpace(hl[idx+1:])
		header := Header{Key: partOne, Value: partTwo}
		req.Headers = append(req.Headers, header)
	}

	return req, nil
}
