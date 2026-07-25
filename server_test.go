package main

import (
	"io"
	"net"
	"testing"
)

func TestHandleConnection(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()

	// must be on a separate worker, otherwise the actual logic will block infinitely.
	go handleConnection(server)

	_, err := client.Write([]byte("GET / HTTP/1.1\r\n"))
	if err != nil {
		t.Fatalf("error writing mock request: %v", err)
	}

	got, err := io.ReadAll(client)
	if err != nil {
		t.Fatalf("error reading response from server: %v", err)
	}

	want := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello, World!"

	if string(got) != want {
		t.Errorf("got %q, want %q", string(got), want)
	}
}
