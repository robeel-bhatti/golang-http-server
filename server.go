package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Serve(network string, addr string) {

	// create a socket and bind it to port 8080 on the host machine
	// returns a listener that produces new connections per client
	// and also sets up the kernel backlog queue.
	l, err := net.Listen(network, addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	// create a channel that OS signals get relayed to.
	// we specify the specific signals we want to listen for
	// that indicate termination of the web server.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// spawn a separate gouroutine that blocks until it receives
	// a value from the channel.
	// this prevents instant death of the web server and instead allows us
	// to control shutdown logic.
	go func() {
		s := <-c
		log.Printf("caught OS signal: %v. shutting down...", s)
		l.Close()
	}()

	// an infinite loop to keep listening for incoming TCP connections.
	for {
		// accept an incoming TCP connection from a client
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Print("listener shutdown")
				break
			}
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		// spawn a goroutine for each client connection
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// create a reader that persists across the lifecycle of a TCP connection.
	// this ensures we store bytes that aren't on the underlying socket anymore.
	// if we created a new reader on every iteration, we could potentially discard unused bytes
	// that weren't on the underlying socket and thus are lost forever.
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	router := NewRouter()
	defer conn.Close()

	// infinite loop which allows us to keep reading bytes
	// off of the same TCP connection (HTTP keep-alive functionality).
	for {
		req, err := ParseRequest(r)
		if err != nil {
			// no more bytes to read off of the socket
			if errors.Is(err, io.EOF) {
				return
			}
			log.Printf("parse error from remote %v: %v", conn.RemoteAddr(), err)
			return
		}
		b, _ := json.MarshalIndent(req, "", "  ")
		fmt.Println(string(b))

		res := router.Route(req)
		c, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(c))

		if err := writeResponse(w, res); err != nil {
			return
		}
		if err := w.Flush(); err != nil {
			return
		}
	}
}

// writeResponse writes the response to the client
// format is:
// [HTTP VERSION] [STATUS CODE] [REASON PHRASE]\r\n
// [HEADER KEY]: [HEADER VALUE]\r\n
// \r\n
// [BODY]
// Note: Body is optional
func writeResponse(w io.Writer, res *Response) error {
	fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", res.Status, http.StatusText(res.Status))

	fmt.Fprintf(w, "Content-Length: %d\r\n", len(res.Body))
	for k, v := range res.Headers {
		fmt.Fprintf(w, "%s: %s\r\n", k, v)
	}

	_, err := w.Write([]byte("\r\n"))
	_, err = w.Write(res.Body)
	return err
}
