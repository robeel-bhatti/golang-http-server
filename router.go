package main

type Handler func(request *Request) *Response

type RouteMap map[string]map[Method]Handler

type Method string

const (
	Get  Method = "GET"
	Post Method = "POST"
)

type Router struct {
	routes RouteMap
}

type Response struct {
	Status  int
	Headers map[string]string
	Body    []byte
}

// NewRouter returns a new instance of a Router struct
// with an empty, initialized map.
func NewRouter() *Router {
	return &Router{
		routes: getRouteMapping(),
	}
}

func (r *Router) Route(req *Request) *Response {
	res, ok := r.routes[req.Path]
	if !ok {
		return &Response{
			Status: 404,
			Body:   []byte("resource not found"),
		}
	}

	h, ok := res[Method(req.Method)]
	if !ok {
		return &Response{
			Status: 405,
			Body:   []byte("method not allowed for this resource"),
		}
	}

	return h(req)
}

func getRouteMapping() RouteMap {
	return RouteMap{
		"/ping": {
			Get: Ping,
		},
		"/teapot": {
			Get: Teapot,
		},
	}
}
