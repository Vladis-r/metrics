package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type ConfigServer struct {
	Addr            string `doc:"ip addr for server."`
	StoreInterval   int    `doc:"Interval for save metrics in file."`
	FileStoragePath string `doc:"Path and filename for save metrics."`
	IsRestore       bool   `doc:"Sign that get metrics from file while start server."`
}

type ConfigAgent struct {
	Addr           string `doc:"Ip addr for send metrcis to server."`
	ReportInterval int    `doc:"Interval for send metrics to server."`
	PollInterval   int    `doc:"Interval for update metrics."`
}

// GetConfigServer - config for server.
func GetConfigServer(logger *zap.Logger) *ConfigServer {
	c := &ConfigServer{}
	flag.StringVar(&c.Addr, "a", "localhost:8080", "Address to listen on.")
	flag.IntVar(&c.StoreInterval, "i", 10, "Interval for save metrics in file.")
	flag.StringVar(&c.FileStoragePath, "f", "metric_log.json", "Path to file where save metrics.")
	flag.BoolVar(&c.IsRestore, "r", false, "If true load saved metrics from file while start server.")
	flag.Parse()

	address := strings.ToLower(os.Getenv("ADDRESS"))                   // ip address for server
	storeInterval := strings.ToLower(os.Getenv("STORE_INTERVAL"))      // interval for save metrics in file
	fileStoragePath := strings.ToLower(os.Getenv("FILE_STORAGE_PATH")) // path to file where save metrics
	isRestore := strings.ToLower(os.Getenv("RESTORE"))                 // if true load saved metrics from file while start server

	switch {
	case address != "":
		c.Addr = address
	case storeInterval != "":
		if i, err := strconv.Atoi(storeInterval); err == nil && i > 0 {
			c.StoreInterval = i
		} else {
			logger.Error("invalid STORE_INTERVAL, must be positive: ", zap.String("STORE_INTERVAL", storeInterval))
		}
	case fileStoragePath != "":
		if err := isValidStoragePath(fileStoragePath); err == nil {
			c.FileStoragePath = fileStoragePath
		} else {
			logger.Error("invalid FILE_STORAGE_PATH: ", zap.String("FILE_STORAGE_PATH", fileStoragePath))
		}
	case isRestore != "":
		c.IsRestore = strings.ToLower(isRestore) == "true"
	}

	return c
}

// GetConfigAgent - config for agent.
func GetConfigAgent(logger *zap.Logger) *ConfigAgent {
	c := &ConfigAgent{}

	flag.StringVar(&c.Addr, "a", "localhost:8080", "Address to listen on.")
	flag.IntVar(&c.ReportInterval, "r", 10, "Interval for reporting metrics.")
	flag.IntVar(&c.PollInterval, "p", 2, "Interval for polling metrics.")
	flag.Parse()

	address := strings.ToLower(os.Getenv("ADDRESS"))                // ip address for server
	reportInterval := strings.ToLower(os.Getenv("REPORT_INTERVAL")) // interval for send metrics to server
	poolInterval := strings.ToLower(os.Getenv("POLL_INTERVAL"))     // interval for update metrics

	switch {
	case address != "":
		c.Addr = address
	case reportInterval != "":
		if i, err := strconv.Atoi(reportInterval); err == nil && i > 0 {
			c.ReportInterval = i
		} else {
			logger.Error("invalid REPORT_INTERVAL, must be positive: ", zap.String("REPORT_INTERVAL", reportInterval))
		}
	case poolInterval != "":
		if i, err := strconv.Atoi(poolInterval); err == nil && i > 0 {
			c.PollInterval = i
		} else {
			logger.Error("invalid POLL_INTERVAL, must be positive: ", zap.String("POLL_INTERVAL", poolInterval))
		}
	}

	return c
}

func isValidStoragePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path is not exist: %s", path)
		}
		return fmt.Errorf("path error: %v", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not directory: %s", path)
	}
	if strings.Contains(filepath.Clean(path), "..") {
		return fmt.Errorf("incorrected path with '..'")
	}

	return nil
}
