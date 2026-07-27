package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http/httputil"
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

func (p *Parser) ParseRequest() (*Request, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}

	method, path, version, err := parseRequestLine(line)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	for {
		headerLine, err := p.readLine()
		if err != nil {
			return nil, err
		}
		if headerLine == "\r\n" {
			break
		}
		key, value, err := parseHeaders(line)
		if err != nil {
			return nil, err
		}
		headers[key] = value
	}

	te, _ := headers["transfer-encoding"]
	cl, _ := headers["content-length"]
	body, err := p.parseBody(te, cl)
	return newRequest(method, path, version, headers, body), nil
}

func (p *Parser) readLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line, nil
}

func parseRequestLine(line string) (string, string, string, error) {
	reqSlice := strings.Split(strings.TrimSpace(line), " ")
	if len(reqSlice) != 3 {
		return "", "", "", fmt.Errorf("invalid request line: %q", line)
	}
	return reqSlice[0], reqSlice[1], reqSlice[2], nil
}

func parseHeaders(line string) (string, string, error) {
	headerSlice := strings.Split(strings.TrimSpace(line), ":")
	if len(headerSlice) != 2 {
		return "", "", fmt.Errorf("invalid header line: %q", line)
	}
	headerKey := strings.ToLower(strings.TrimSpace(headerSlice[0]))
	headerValue := strings.TrimSpace(headerSlice[1])
	return headerKey, headerValue, nil
}

func (p *Parser) parseBody(te, cl string) ([]byte, error) {
	var r io.Reader
	if strings.Contains(te, "chunked") {
		r = httputil.NewChunkedReader(p.reader)
	} else if cl != "" {
		clNum, _ := strconv.Atoi(cl)
		r = io.LimitReader(p.reader, int64(clNum))
	} else {
		return nil, fmt.Errorf("no content-length or transfer-encoding specified")
	}

	body, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return body, nil
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
