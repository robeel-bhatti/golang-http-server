package main

const (
	Network = "tcp"
	Address = ":8080"
)

func main() {
	Serve(Network, Address)
}
