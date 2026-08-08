package main

import (
	"bufio"
	"fmt"
	"io"
)

type ResponseWriter struct {
	writer *bufio.Writer
}

type Response struct {
	Status  int
	Headers map[string]string
	Body    []byte
}

func NewResponseWriter(w io.Writer) *ResponseWriter {
	return &ResponseWriter{
		writer: bufio.NewWriter(w),
	}
}

func (w *ResponseWriter) Write(resp *Response) error {
	fmt.Fprintf(w.writer, "HTTP/1.1 %d %s\r\n", resp.Status, resolveStatus(resp.Status))
	for k, v := range resp.Headers {
		fmt.Fprintf(w.writer, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(w.writer, "Content-Length: %d\r\n", len(resp.Body))
	w.writer.WriteString("\r\n")
	w.writer.Write(resp.Body)
	return w.writer.Flush()
}

func resolveStatus(s int) string {
	switch s {
	case 200:
		return "OK"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
}
