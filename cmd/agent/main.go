package main

import (
	"log"

	"github.com/Vladis-r/metrics.git/cmd/config"
	"github.com/Vladis-r/metrics.git/internal/agent"
	"github.com/Vladis-r/metrics.git/internal/middleware"
	models "github.com/Vladis-r/metrics.git/internal/model"
	"go.uber.org/zap"
)

func main() {
	logger, err := middleware.InitLogger() // create logger
	if err != nil {
		log.Fatalf("Cant create logger: %v", err)
	}
	defer logger.Sync()
	cfg := config.GetConfigAgent(logger) // Parse command-line arguments.
	m := models.NewMetricsMap()          // Init client and map for metrics.

	logger.Info("Start metrics agent...")
	logger.Info("With config:\n %v", zap.Any("Config:", cfg))

	goroutines := []func(*models.MetricsMap, *config.ConfigAgent){
		agent.GoUpdateMetrics,
		agent.GoReportMetrics,
	}

	for _, goroutine := range goroutines {
		m.Wg.Add(1)
		go goroutine(m, cfg)
	}

	logger.Info("Press Ctrl+C to exit")

	m.Wg.Wait()
}
