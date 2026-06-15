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
	for {
		rline, err := r.ReadString('\n')
		if err != nil {
			log.Printf("failed to read: %v", err)
			break
		}
		rline = strings.TrimRight(rline, "\r\n")
		log.Print(rline)

		for {
			hline, err := r.ReadString('\n')
			if err != nil {

			}

			hline = strings.TrimRight(hline, "\r\n")
			if hline == "" {
				return
			}
			log.Print(hline)
		}
	}
}
