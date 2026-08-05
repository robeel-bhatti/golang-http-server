package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var IdleTimeout = 5 * time.Second

func Serve(protocol, port string) {
	ln, err := net.Listen(protocol, port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serve(ln)
}

func serve(ln net.Listener) {
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("new connection from %v", conn.RemoteAddr())

	parser := NewParser(conn)
	writer := NewResponseWriter(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(IdleTimeout))

		keepAlive, err := serveRequest(parser, writer)
		if err != nil {
			logConnError(conn, err)
			return
		}
		if !keepAlive {
			return
		}
	}
}

func serveRequest(p *Parser, w *ResponseWriter) (keepAlive bool, err error) {
	req, err := p.ParseRequest()
	if err != nil {
		return false, err
	}

	mockRes := getMockResponse()
	keepAlive = !strings.EqualFold(req.Headers["connection"], "close")
	if !keepAlive {
		mockRes.Headers["Connection"] = "close"
	}
	return keepAlive, w.Write(mockRes)
}

func logConnError(conn net.Conn, err error) {
	switch {
	case errors.Is(err, io.EOF):
	case errors.Is(err, os.ErrDeadlineExceeded):
		log.Printf("deadline exceeded for peer %v", conn.RemoteAddr())
	default:
		log.Printf("unexpected error for peer %v: %v", conn.RemoteAddr(), err)
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
