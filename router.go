package main

type Handler func(*Request) *Response

type ServeMux struct {
	routes map[string]map[string]Handler
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		routes: make(map[string]map[string]Handler),
	}
}

func (m *ServeMux) Handle(pattern, method string, handler Handler) {
	innerMap := make(map[string]Handler)
	innerMap[method] = handler
	m.routes[pattern] = innerMap
}

func (m *ServeMux) Dispatch(r *Request) *Response {
	res, ok := m.routes[r.Metadata.Path]
	if !ok {
		return &Response{Status: 404, Body: []byte("resource not found")}
	}
	h, ok := res[r.Metadata.Method]
	if !ok {
		return &Response{Status: 405, Body: []byte("method not allowed")}
	}
	return h(r)
}

func GetHello(r *Request) *Response {
	return &Response{
		Status:  200,
		Headers: map[string]string{"Content-Type": "text/plain"},
		Body:    []byte("Hello, World!"),
	}
}
