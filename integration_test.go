package main

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"testing"
)

const (
	TestProtocol = "tcp"
	TestPort     = ":8081"
)

func TestServe_Integration(t *testing.T) {
	serverReady := make(chan struct{})
	go func() {
		close(serverReady)
		Serve(TestProtocol, TestPort)
	}()

	<-serverReady

	conn, err := net.Dial(TestProtocol, TestPort)
	if err != nil {
		t.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	httpRequest := "GET /health HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"Connection: close\r\n" +
		"\r\n"
	_, err = conn.Write([]byte(httpRequest))
	if err != nil {
		t.Fatalf("could not write to server: %v", err)
	}

	expectedBody := "Hello, World!"
	expectedContentLength := int64(len(expectedBody))

	reader := bufio.NewReader(conn)
	res, err := http.ReadResponse(reader, nil)

	if err != nil {
		t.Fatalf("could not read response from server: %v", err)
	}
	if res.StatusCode != 200 {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}
	if res.Close != true {
		t.Errorf("expected close to be true, got %t", res.Close)
	}
	if res.ContentLength != expectedContentLength {
		t.Errorf("expected content length to be %d, got %d", expectedContentLength, res.ContentLength)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read body from server response: %v", err)
	}
	if string(bodyBytes) != expectedBody {
		t.Errorf("expected body to be %q, got %q", expectedBody, string(bodyBytes))
	}
}
