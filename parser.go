package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

type Parser struct {
	reader *bufio.Reader
}

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
}

// validMethods is a set of HTTP methods the server accepts
var validMethods = []string{
	"GET",
	"POST",
	"PUT",
	"DELETE",
	"PATCH",
	"HEAD",
	"OPTIONS",
}

func NewParser(r io.Reader) *Parser {
	return &Parser{reader: bufio.NewReader(r)}
}

func (p *Parser) ParseRequest() (*Request, error) {
	method, path, version, err := p.parseRequestLine()
	if err != nil {
		return nil, err
	}

	headers, err := p.parseHeaders()
	if err != nil {
		return nil, err
	}

	body, err := p.parseBody(headers["content-length"])
	if err != nil {
		return nil, err
	}

	return &Request{
		Method:  method,
		Path:    path,
		Version: version,
		Headers: headers,
		Body:    body,
	}, nil
}

func (p *Parser) parseRequestLine() (method, path, version string, err error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}

	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid request line: %s", line)
	}
	if !slices.Contains(validMethods, parts[0]) {
		return "", "", "", fmt.Errorf("invalid HTTP verb: %s", parts[0])
	}

	return parts[0], parts[1], parts[2], nil
}

func (p *Parser) parseHeaders() (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := p.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			return headers, nil
		}

		key, value, found := strings.Cut(line, ":")
		if !found {
			return nil, fmt.Errorf("invalid HTTP header: %q", line)
		}
		headers[strings.ToLower(strings.TrimSpace(key))] = strings.TrimSpace(value)
	}
}

func (p *Parser) parseBody(contentLength string) ([]byte, error) {
	if contentLength == "" {
		return nil, nil
	}

	n, err := strconv.Atoi(contentLength)
	if err != nil || n < 0 {
		return nil, fmt.Errorf("invalid content length: %q", contentLength)
	}
	return io.ReadAll(io.LimitReader(p.reader, int64(n)))
}
