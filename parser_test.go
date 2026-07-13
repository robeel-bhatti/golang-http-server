package main

import (
	"bufio"
	"strings"
	"testing"
)

type TestRequest struct {
	name       string
	raw        string
	wantMethod string
	wantPath   string
}

func TestParseRequest(t *testing.T) {
	tests := []TestRequest{
		{
			name:       "success_parse_ping_request",
			raw:        "GET /ping HTTP/1.1\r\nHost: example.com\r\n\r\n",
			wantMethod: "GET",
			wantPath:   "/ping",
		},
		{
			name:       "success_parse_teapot_request",
			raw:        "GET /teapot HTTP/1.1\r\nHost: example.com\r\n\r\n",
			wantMethod: "GET",
			wantPath:   "/teapot",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(tt.raw))
			req, err := ParseRequest(r)
			if err != nil {
				t.Fatalf("ParseRequest returned error: %v", err)
			}
			if req.Method != tt.wantMethod {
				t.Errorf("ParseRequest method = %v, want %v", req.Method, tt.wantMethod)
			}
			if req.Path != tt.wantPath {
				t.Errorf("ParseRequest path = %v, want %v", req.Path, tt.wantPath)
			}
		})
	}
}
