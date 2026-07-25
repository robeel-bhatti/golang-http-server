package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	line, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("failed to read: %v", err)
		return
	}

	reqLine := strings.TrimSpace(line)
	reqParts := strings.Split(reqLine, " ")

	method := reqParts[0]
	path := reqParts[1]
	version := reqParts[2]

	log.Printf("[METHOD] - %s, [PATH] - %s, [VERSION] - %s", method, path, version)
	writeResponse(conn)
}

func writeResponse(conn net.Conn) {
	w := bufio.NewWriter(conn)
	body := "Hello, World!"
	w.WriteString("HTTP/1.1 200 OK\r\n")
	w.WriteString("Content-Type: text/plain\r\n")
	w.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
	w.WriteString("\r\n")
	w.WriteString(body)

	err := w.Flush()
	if err != nil {
		log.Printf("failed to write response to remote %v: %v", conn.RemoteAddr(), err)
	}
}
