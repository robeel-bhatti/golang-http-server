package main

import (
	"bufio"
	"errors"
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
	req, err := readRequest(r)
	if err != nil {
		return nil, err
	}

	headers, err := readHeaders(r)
	if err != nil {
		return nil, err
	}
	
	body, err := readBody(r, req)
	if err != nil {
		return nil, err
	}

	req.Headers = headers
	req.Body = body
	return req, nil
}

func readRequest(r *bufio.Reader) (*Request, error) {
	line, err := readLine(r)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("malformed request: %s", line)
	}
	return &Request{Method: parts[0], Path: parts[1], Version: parts[2]}, nil
}

func readHeaders(r *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := readLine(r)
		if err != nil {
			return nil, err
		}
		if line == "" {
			break
		}

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return nil, fmt.Errorf("malformed header: %s", line)
		}

		key = strings.ToLower(strings.TrimSpace(key))
		headers[key] = strings.TrimSpace(value)

	}
	return headers, nil
}

func readBody(r *bufio.Reader, req *Request) ([]byte, error) {
	cl, ok := req.Headers["content-length"]
	if !ok {
		return nil, nil
	}

	n, err := strconv.Atoi(cl)
	if err != nil || n < 0 {
		return nil, fmt.Errorf("invalid content-length: %v", cl)
	}

	body := make([]byte, n)
	if _, err := io.ReadFull(r, body); err != nil {
		return nil, err
	}

	return body, nil
}

// readLine reads on CRLF/LF-terminated line and strips the terminator
// a partial line cut short by EOF is reported as ErrUnexpectedEOF
func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) && line != "" {
			return "", io.ErrUnexpectedEOF
		}
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}
