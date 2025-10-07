package flags

import (
	"flag"
)

var (
	Addr           string
	ReportInterval int
	PollInterval   int
)

func init() {
	flag.StringVar(&Addr, "a", "localhost:8080", "Address to listen on")
	flag.IntVar(&ReportInterval, "r", 10, "Interval for reporting metrics")
	flag.IntVar(&PollInterval, "p", 2, "Interval for polling metrics")
}
