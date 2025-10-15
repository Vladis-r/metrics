package config

import (
	"flag"
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

	return c
}
