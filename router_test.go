package main

import (
	"maps"
	"testing"
)

func TestServeMux_Dispatch(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		wantStatus  int
		wantHeaders map[string]string
	}{
		{"GET is dispatched", "GET", "/hello", 200, nil},
		{"resource not found", "GET", "/nope", 404, nil},
		{"method not allowed", "POST", "/hello", 405, map[string]string{"Allow": "GET"}},
	}

	m := NewServeMux()
	m.Handle("/hello", "GET", GetHello)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{Metadata: &Metadata{Method: tt.method, Path: tt.path}}
			res := m.Dispatch(req)

			if tt.wantStatus != res.Status {
				t.Errorf("want status %d, got %d", tt.wantStatus, res.Status)
			}
			if tt.wantHeaders != nil && !maps.Equal(tt.wantHeaders, res.Headers) {
				t.Errorf("want headers %v, got %v", tt.wantHeaders, res.Headers)
			}
		})
	}
}
