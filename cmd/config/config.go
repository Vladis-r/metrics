package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Addr           string
	ReportInterval int
	PollInterval   int
}

func GetConfig() *Config {
	c := &Config{}
	flag.StringVar(&c.Addr, "a", "localhost:8080", "Address to listen on")
	flag.IntVar(&c.ReportInterval, "r", 10, "Interval for reporting metrics")
	flag.IntVar(&c.PollInterval, "p", 2, "Interval for polling metrics")

	flag.Parse()

	switch {
	case os.Getenv("ADDR") != "":
		c.Addr = os.Getenv("ADDR")
	case os.Getenv("REPORT_INTERVAL") != "":
		if i, err := strconv.Atoi(os.Getenv("REPORT_INTERVAL")); err == nil {
			if i > 0 {
				c.ReportInterval = i
			}
		}
	case os.Getenv("POLL_INTERVAL") != "":
		if i, err := strconv.Atoi(os.Getenv("POLL_INTERVAL")); err == nil {
			if i > 0 {
				c.PollInterval = i
			}
		}
	}
	return c
}
