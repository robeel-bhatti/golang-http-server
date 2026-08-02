package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const IdleTimeout = 5 * time.Second

func Serve(protocol, port string) {
	listener, err := net.Listen(protocol, port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	log.Printf("new connection from %v", conn.RemoteAddr())
	defer conn.Close()
	parser := NewParser(conn)
	writer := NewResponseWriter(conn)

	for {
		_ = conn.SetReadDeadline(time.Now().Add(IdleTimeout))

		req, err := parser.ParseRequest()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			if errors.Is(err, os.ErrDeadlineExceeded) {
				log.Printf("deadline exceeded for peer %v", conn.RemoteAddr())
				return
			}
			log.Printf("failed to parse request: %v", err)
			return
		}

		mockRes := getMockResponse()
		shouldClose := req.Headers["connection"] == "close"
		if shouldClose {
			mockRes.Headers["Connection"] = "close"
		}

		err = writer.Write(mockRes)
		if err != nil {
			log.Printf("failed to write response to peer %v: %v", conn.RemoteAddr(), err)
			return
		}
		if shouldClose {
			return
		}
	}
}

func getMockResponse() *Response {
	body := "Hello, World!"
	return &Response{
		Status: 200,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		Body: []byte(body),
	}
}
