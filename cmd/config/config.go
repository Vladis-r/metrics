package config

import (
	"flag"
	"fmt"
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

	if addr := os.Getenv("ADDR"); addr != "" {
		c.Addr = addr
	}
	if reportStr := os.Getenv("REPORT_INTERVAL"); reportStr != "" {
		if i, err := strconv.Atoi(reportStr); err == nil && i > 0 {
			c.ReportInterval = i
		} else {
			fmt.Printf("Invalid REPORT_INTERVAL, must be positive: %s\n", reportStr)
		}
	}
	if pollStr := os.Getenv("POLL_INTERVAL"); pollStr != "" {
		if i, err := strconv.Atoi(pollStr); err == nil && i > 0 {
			c.PollInterval = i
		} else {
			fmt.Printf("Invalid POLL_INTERVAL, must be positive: %s\n", pollStr)
		}
	}
	return c
}
