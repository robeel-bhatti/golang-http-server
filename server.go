package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Serve(network string, addr string) {
	// first open up a socket endpoint on the provided port
	// that listens for incoming TCP connections
	l, err := net.Listen(network, addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	defer l.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-c
		log.Printf("caught OS signal: %v. shutting down...", s)
		l.Close()
	}()

	// constantly listening...
	for {
		// accept an incoming TCP connection from a client
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Printf("listener shutdown")
				return
			}
			log.Printf("failed to accept connection from remote %v: %v", conn.RemoteAddr(), err)
		}

		// spawn a goroutine for each client connection
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	r := bufio.NewReader(conn)
	req, err := ParseRequest(r)
}
