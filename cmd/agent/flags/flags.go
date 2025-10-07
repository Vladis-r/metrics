package flags

import (
	"flag"
	"time"
)

var (
	Addr           string
	ReportInterval time.Duration
	PollInterval   time.Duration
)

func init() {
	flag.StringVar(&Addr, "a", "localhost:8080", "Address to listen on")
	flag.DurationVar(&ReportInterval, "r", 10*time.Second, "Interval for reporting metrics")
	flag.DurationVar(&PollInterval, "p", 2*time.Second, "Interval for polling metrics")
}
