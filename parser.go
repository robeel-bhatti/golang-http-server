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

var ValidMethods = []string{
	"GET",
	"POST",
	"PUT",
	"DELETE",
	"PATCH",
	"HEAD",
	"OPTIONS",
}

func (p *Parser) ParseRequest() (*Request, error) {
	method, path, version, err := p.getRequestLineComponents()
	if err != nil {
		return nil, err
	}

	headers, err := p.getHeaders()
	if err != nil {
		return nil, err
	}

	body, err := p.getBody(headers["content-length"])
	if err != nil {
		return nil, err
	}

	return newRequest(method, path, version, headers, body), nil
}

func (p *Parser) getRequestLineComponents() (string, string, string, error) {
	line, err := p.readLine()
	if err != nil {
		return "", "", "", err
	}
	req := strings.Split(strings.TrimSpace(line), " ")
	if len(req) != 3 {
		return "", "", "", fmt.Errorf("invalid request line: %q", line)
	}
	if !slices.Contains(ValidMethods, req[0]) {
		return "", "", "", fmt.Errorf("invalid HTTP method: %q", req[0])
	}
	return req[0], req[1], req[2], nil
}

func (p *Parser) getHeaders() (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := p.readLine()
		if err != nil {
			return nil, err
		}
		l := strings.TrimSpace(line)
		if l == "" {
			break
		}
		if len(l) != 2 {
			return nil, fmt.Errorf("invalid header line: %q", line)
		}
		headerList := strings.Split(l, ":")
		key := strings.ToLower(strings.TrimSpace(headerList[0]))
		value := strings.TrimSpace(headerList[1])
		headers[key] = value
	}
	return headers, nil
}

func (p *Parser) getBody(cl string) ([]byte, error) {
	var r io.Reader

	clNum, _ := strconv.Atoi(cl)
	r = io.LimitReader(p.reader, int64(clNum))

	body, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (p *Parser) readLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line, nil
}

func newRequest(method, path, version string, headers map[string]string, body []byte) *Request {
	return &Request{
		Method:  method,
		Path:    path,
		Version: version,
		Headers: headers,
		Body:    body,
	}
}
