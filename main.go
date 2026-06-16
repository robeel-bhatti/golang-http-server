package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Request struct {
	Method  string
	Path    string
	Headers []Header
}

type Header struct {
	Key   string
	Value string
}

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	sigCh := make(chan os.Signal, 1) // buffered channel otherwise sender will block and discard signal
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// spawn a separate goroutine that blocks and listens to channel
	go func() {
		sig := <-sigCh
		log.Printf("received signal: %v. shutting down...", sig)
		l.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("listener shut down")
				break
			}
			log.Printf("failed to accept: %v", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	r := bufio.NewReader(conn)
	defer conn.Close()
	for {
		requestLine, err := r.ReadString('\n')
		if err != nil {
			log.Printf("failed to read: %v", err)
			return
		}
		req := parseRequestLine(requestLine)
		var headers []Header

		for {
			headerLine, err := r.ReadString('\n')
			if err != nil {
				log.Printf("failed to read: %v", err)
				return
			}
			headerLine = strings.TrimRight(headerLine, "\r\n")
			if headerLine == "" {
				break
			}
			headers = append(headers, parseHeader(headerLine))
		}
		req.Headers = headers
	}
}

func parseRequestLine(requestLine string) *Request {
	requestLine = strings.TrimRight(requestLine, "\r\n")
	parts := strings.Split(requestLine, " ")
	return &Request{Method: parts[0], Path: parts[1]}
}

func parseHeader(headerLine string) Header {
	parts := strings.Split(headerLine, ": ")
	return Header{Key: parts[0], Value: parts[1]}
}
