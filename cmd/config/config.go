package config

import (
	"flag"
	"net"
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
	var (
		address         string
		storeInterval   string
		fileStoragePath string
		isRestore       string
	)
	c := &ConfigServer{}
	flag.StringVar(&address, "a", "localhost:8080", "Address to listen on.")
	flag.StringVar(&storeInterval, "i", "300", "Interval for save metrics in file.")
	flag.StringVar(&fileStoragePath, "f", "metric_log.json", "Path to file where save metrics.")
	flag.StringVar(&isRestore, "r", "true", "If true load saved metrics from file while start server.")
	flag.Parse()

	address = getEnv("ADDRESS", address)                           // ip address for server
	storeInterval = getEnv("STORE_INTERVAL", storeInterval)        // interval for save metrics in file
	fileStoragePath = getEnv("FILE_STORAGE_PATH", fileStoragePath) // path to file where save metrics
	isRestore = getEnv("RESTORE", isRestore)                       // if true load saved metrics from file while start server

	c.Addr = validAddress(address, logger)
	c.StoreInterval = validStoreInterval(storeInterval, logger)
	c.FileStoragePath = validStoragePath(fileStoragePath, logger)
	c.IsRestore = validRestore(isRestore, logger)

	return c
}

// GetConfigAgent - config for agent.
func GetConfigAgent(logger *zap.Logger) *ConfigAgent {
	var (
		address        string
		reportInterval string
		pollInterval   string
	)
	c := &ConfigAgent{}

	flag.StringVar(&address, "a", "localhost:8080", "Address to listen on.")
	flag.StringVar(&reportInterval, "r", "10", "Interval for reporting metrics.")
	flag.StringVar(&pollInterval, "p", "2", "Interval for polling metrics.")
	flag.Parse()

	address = getEnv("ADDRESS", address)                       // ip address for server
	reportInterval = getEnv("REPORT_INTERVAL", reportInterval) // interval for send metrics to server
	pollInterval = getEnv("POLL_INTERVAL", pollInterval)       // interval for update metrics

	c.Addr = validAddress(address, logger)
	c.ReportInterval = validReportInterval(reportInterval, logger)
	c.PollInterval = validPollInterval(pollInterval, logger)

	return c
}

func validAddress(ipAddr string, logger *zap.Logger) string {
	host, _, err := net.SplitHostPort(ipAddr)
	if err != nil {
		logger.Error("invalid ADDRESS while splitHostPort: ", zap.String("ADDRESS", ipAddr), zap.String("Host", host))
		return ""
	}

	parsed := net.ParseIP(host)
	if parsed == nil && host != "localhost" {
		logger.Error("invalid ADDRESS while ParseIP: ", zap.String("ADDRESS", ipAddr), zap.String("Host", host))
		return ""
	}
	return ipAddr
}

func validPollInterval(pollInterval string, logger *zap.Logger) int {
	if interval, err := strconv.Atoi(pollInterval); err == nil && interval > 0 {
		return interval
	} else {
		logger.Error("invalid POLL_INTERVAL, must be positive: ", zap.String("POLL_INTERVAL", pollInterval))
	}
	return 0
}

func validReportInterval(reportInterval string, logger *zap.Logger) int {
	if interval, err := strconv.Atoi(reportInterval); err == nil && interval > 0 {
		return interval
	} else {
		logger.Error("invalid REPORT_INTERVAL, must be positive: ", zap.String("REPORT_INTERVAL", reportInterval))
	}
	return 0
}

func validRestore(isRestore string, logger *zap.Logger) bool {
	restore := strings.ToLower(isRestore)
	if restore == "true" {
		return true
	}
	if restore != "false" {
		logger.Error("invalid RESTORE, must be true or false: ", zap.String("RESTORE", isRestore))
	}
	return false
}

func validStoreInterval(storeInterval string, logger *zap.Logger) int {
	if interval, err := strconv.Atoi(storeInterval); err == nil && interval > 0 {
		return interval
	} else {
		logger.Error("invalid STORE_INTERVAL, must be positive: ", zap.String("STORE_INTERVAL", storeInterval))
	}
	return 0
}

func validStoragePath(path string, logger *zap.Logger) string {
	cleanPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		logger.Error("incorrected path", zap.String("Path", path))
		return ""
	}

	if strings.Contains(cleanPath, "..") {
		logger.Error("incorrected path with '..'", zap.String("Path", path))
		return ""
	}

	return cleanPath
}

// getEnv - return value from ENV by key or default.
func getEnv(key, defaultValue string) (value string) {
	if os.Getenv(key) != "" {
		return value
	}
	return defaultValue
}
