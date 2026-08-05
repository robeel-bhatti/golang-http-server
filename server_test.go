package main

import (
	"bytes"
	"strings"
	"testing"
)

// runServeRequest returns out to verify w.Write() behavior
func runServeRequest(t *testing.T, raw string) (keepAlive bool, out string, err error) {
	t.Helper()
	p := NewParser(strings.NewReader(raw))
	var buf bytes.Buffer
	w := NewResponseWriter(&buf)
	keepAlive, err = serveRequest(p, w)
	return keepAlive, buf.String(), err
}

func TestServeRequest_ValidGET(t *testing.T) {
	keepAlive, out, err := runServeRequest(t, "GET / HTTP/1.1\r\nHost:x\r\n\r\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !keepAlive {
		t.Error("expected keepAlive to be true")
	}
	if !strings.Contains(out, "HTTP/1.1 200 OK") {
		t.Error("missing status line, got :\n%", out)
	}
	if !strings.Contains(out, "Hello, World!") {
		t.Error("missing body, got :\n%", out)
	}

}
