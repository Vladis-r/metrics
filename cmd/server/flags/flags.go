package flags

import "flag"

var (
	Addr string
)

func init() {
	flag.StringVar(&Addr, "a", "localhost:8080", "Address to listen on")
}
