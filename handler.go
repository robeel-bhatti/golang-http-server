package main

// Ping validates a body can be returned
func Ping(request *Request) *Response {
	return &Response{
		Status:  200,
		Headers: map[string]string{"Content-Type": "text/plain"},
		Body:    []byte("pong"),
	}
}

// Teapot validates an arbitrary status code can be returned
func Teapot(request *Request) *Response {
	return &Response{Status: 418}
}
