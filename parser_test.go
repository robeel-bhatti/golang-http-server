package main

import (
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
			name:  "Bodyless GET",
			input: "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n",
			want: &Request{
				Method:  "GET",
				Path:    "/",
				Version: "HTTP/1.1",
				Headers: map[string]string{
					"host": "example.com",
				},
				Body: nil,
			},
		},
		{
			name:  "POST with a body",
			input: "POST / HTTP/1.1\r\nContent-Length: 12\r\n\r\nHello World!",
			want: &Request{
				Method:  "POST",
				Path:    "/",
				Version: "HTTP/1.1",
				Headers: map[string]string{
					"content-length": "12",
				},
				Body: []byte("Hello World!"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewParser(strings.NewReader(tt.input)).ParseRequest()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot %+v\nwant %+v", got, tt.want)
			}
		})
	}
}
