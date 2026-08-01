package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrMalformedRequest = errors.New("malformed request")
	ErrBadRequest       = errors.New("bad request")
)

type Parser struct {
	reader *bufio.Reader
}

type Metadata struct {
	Method  string
	Path    string
	Version string
}

type Request struct {
	Metadata *Metadata
	Headers  map[string]string
	Body     []byte
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
	metadata, err := p.parseRequestLine()
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
		Metadata: metadata,
		Headers:  headers,
		Body:     body,
	}, nil
}

func (p *Parser) parseRequestLine() (metadata *Metadata, err error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			if line == "" {
				return nil, io.EOF
			}
			return nil, io.ErrUnexpectedEOF
		}
		return nil, err
	}

	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("%w. invalid request line: %s", ErrMalformedRequest, line)
	}
	if !slices.Contains(validMethods, parts[0]) {
		return nil, fmt.Errorf("%w. invalid HTTP verb: %s", ErrBadRequest, parts[0])
	}

	return &Metadata{
		Method:  parts[0],
		Path:    parts[1],
		Version: parts[2],
	}, nil
}

func (p *Parser) parseHeaders() (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := p.reader.ReadString('\n')
		if err != nil {
			return nil, isUnexpectedEOF(err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			return headers, nil
		}

		key, value, found := strings.Cut(line, ":")
		if !found {
			return nil, fmt.Errorf("%w. invalid HTTP header: %q", ErrMalformedRequest, line)
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
		return nil, fmt.Errorf("%w. invalid content length: %q", ErrBadRequest, contentLength)
	}

	if n == 0 {
		return nil, nil
	}

	buf := make([]byte, n)
	_, err = io.ReadFull(p.reader, buf)
	if err != nil {
		return nil, isUnexpectedEOF(err)
	}
	return buf, nil
}

func isUnexpectedEOF(err error) error {
	if errors.Is(err, io.EOF) {
		return io.ErrUnexpectedEOF
	}
	return err
}
