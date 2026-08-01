package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

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
	reader := bufio.NewReader(conn)
	parser := NewParser(reader)

	for {
		_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		req, err := parser.ParseRequest()
		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				log.Printf("connection from %v closed unexpectedly", conn.RemoteAddr())
				return
			}
			log.Printf("failed to parse request: %v", err)
			return
		}

		cc := false
		_, ok := req.Headers["connection"]
		if ok && req.Headers["connection"] == "close" {
			cc = true
		}
		writeResponse(conn, cc)

		if cc {
			return
		}
	}
}

func writeResponse(conn net.Conn, cc bool) {
	w := bufio.NewWriter(conn)
	body := "Hello, World!"
	w.WriteString("HTTP/1.1 200 OK\r\n")
	w.WriteString("Content-Type: text/plain\r\n")
	w.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
	if cc {
		w.WriteString("Connection: close\r\n")
	}
	w.WriteString("\r\n")
	w.WriteString(body)

	err := w.Flush()
	if err != nil {
		log.Printf("failed to write response to remote %v: %v", conn.RemoteAddr(), err)
	}
}
