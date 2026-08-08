package main

import "strings"

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
	if m.routes[pattern] == nil {
		m.routes[pattern] = make(map[string]Handler)
	}
	m.routes[pattern][method] = handler
}

func (m *ServeMux) Dispatch(r *Request) *Response {
	methods, ok := m.routes[r.Metadata.Path]
	if !ok {
		return &Response{Status: 404, Body: []byte("resource not found")}
	}
	h, ok := methods[r.Metadata.Method]
	if !ok {
		var allowed []string
		for key := range methods {
			allowed = append(allowed, key)
		}
		return &Response{
			Status: 405,
			Headers: map[string]string{
				"Allow": strings.Join(allowed, ", "),
			},
			Body: []byte("method not allowed"),
		}
	}
	return h(r)
}
