package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

// Response represents an HTTP response produced by a handler,
// ready to be serialized onto the wire.
type Response struct {
	Status  int
	Headers map[string]string
	Body    []byte
}

// Write serializes the response onto w in HTTP/1.1 wire format:
//
//	HTTP/1.1 [STATUS CODE] [REASON PHRASE]\r\n
//	[HEADER KEY]: [HEADER VALUE]\r\n
//	\r\n
//	[BODY]
//
// It flushes the writer before returning, so a nil error means the
// full response reached the underlying connection.
func (r *Response) Write(w *bufio.Writer) error {
	fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", r.Status, http.StatusText(r.Status))

	// Content-Length is derived from the body, so this method owns it.
	// Skip any handler-set copy to avoid emitting two conflicting values.
	fmt.Fprintf(w, "Content-Length: %d\r\n", len(r.Body))

	for k, v := range r.Headers {
		if strings.EqualFold(k, "Content-Length") {
			continue
		}
		fmt.Fprintf(w, "%s: %s\r\n", k, v)
	}

	// blank line separates the headers from the body
	w.WriteString("\r\n")
	w.Write(r.Body)

	// bufio.Writer latches the first error any write above hit and
	// returns it from Flush, so this single check covers all of them.
	return w.Flush()
}
