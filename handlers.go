package main

func GetHello(r *Request) *Response {
	return &Response{
		Status:  200,
		Headers: map[string]string{"Content-Type": "text/plain"},
		Body:    []byte("Hello, World!"),
	}
}
