package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestParser_ParseRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Request
		wantErr error
	}{
		{
			name:  "GET with no body",
			input: "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n",
			want: &Request{
				Metadata: &Metadata{
					Method:  "GET",
					Path:    "/",
					Version: "HTTP/1.1",
				},
				Headers: map[string]string{
					"host": "example.com",
				},
				Body: nil,
			},
			wantErr: nil,
		},
		{
			name:    "GET with malformed request line",
			input:   "GET invalidText / HTTP/1.1\r\nHost: example.com\r\n\r\n",
			want:    nil,
			wantErr: ErrMalformedRequest,
		},
		{
			name:  "POST with a body",
			input: "POST / HTTP/1.1\r\nContent-Length: 12\r\n\r\nHello World!",
			want: &Request{
				Metadata: &Metadata{
					Method:  "POST",
					Path:    "/",
					Version: "HTTP/1.1",
				},
				Headers: map[string]string{
					"content-length": "12",
				},
				Body: []byte("Hello World!"),
			},
			wantErr: nil,
		},
		{
			name:    "POST with invalid formatted header",
			input:   "POST / HTTP/1.1\r\nContent-Length 20\r\n\r\nHello, World!",
			want:    nil,
			wantErr: ErrMalformedRequest,
		},
		{
			name:    "POST with less than 0 content-length header",
			input:   "POST / HTTP/1.1\r\nContent-Length: -2\r\n\r\nHello, World!",
			want:    nil,
			wantErr: ErrBadRequest,
		},
		{
			name:  "POST with empty content-length header",
			input: "POST / HTTP/1.1\r\nContent-Length: \r\n\r\nHello, World!",
			want: &Request{
				Metadata: &Metadata{
					Method:  "POST",
					Path:    "/",
					Version: "HTTP/1.1",
				},
				Headers: map[string]string{
					"content-length": "",
				},
				Body: nil,
			},
			wantErr: nil,
		},
		{
			name:    "Request with unknown verb",
			input:   "GETPOST / HTTP/1.1\r\nContent-Length: \r\n\r\nHello, World!",
			want:    nil,
			wantErr: ErrBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewParser(strings.NewReader(tt.input)).ParseRequest()
			if err != nil && tt.wantErr == nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("\ngot %v\nwant %v", tt.wantErr, err)
			}
			if tt.want != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot %+v\nwant %+v", got, tt.want)
			}
		})
	}
}
