package main

import (
	"errors"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Fatalf("listener has shutdown")
				// end program execution
			}
			log.Printf("failed to accept: %v", err)

			// skip processing this connection and move to the next one
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {

}
