package main

const (
	Protocol = "tcp"
	Port     = ":8080"
)

func main() {
	Serve(Protocol, Port)
}
