package main

import "testing"

func TestGetHello(t *testing.T) {
	m := NewServeMux()
	m.Handle("/hello", "GET", GetHello)
	res := m.Dispatch(&Request{
		Metadata: &Metadata{
			Method:  "GET",
			Path:    "/hello",
			Version: "HTTP/1.1",
		},
		Headers: map[string]string{},
		Body:    nil,
	})

	if res.Status != 200 {
		t.Errorf("expected status 200, got %d", res.Status)
	}
	ct, ok := res.Headers["Content-Type"]
	if !ok {
		t.Errorf("expected Content-Type header to be set")
	} else if ct != "text/plain" {
		t.Errorf("expected Content-Type header to be text/plain, got %s", ct)
	}
	if string(res.Body) != "Hello, World!" {
		t.Errorf("expected body 'Hello, World!', got %s", string(res.Body))
	}
}
