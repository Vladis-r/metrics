package main

import (
	"fmt"

	"github.com/Vladis-r/metrics.git/cmd/config"
	"github.com/Vladis-r/metrics.git/internal/agent"
	models "github.com/Vladis-r/metrics.git/internal/model"
)

func main() {
	cfg := config.GetConfigAgent() // Parse command-line arguments.
	m := models.NewMetricsMap()    // Init client and map for metrics.

	fmt.Println("Start metrics agent...")
	fmt.Printf("With config:\n PollInterval: %v\n ReportInterval: %v\n\n", cfg.PollInterval, cfg.ReportInterval)

	goroutines := []func(*models.MetricsMap, *config.ConfigAgent){
		agent.GoUpdateMetrics,
		agent.GoReportMetrics,
	}

	for _, goroutine := range goroutines {
		m.Wg.Add(1)
		go goroutine(m, cfg)
	}

	fmt.Println("Press Ctrl+C to exit")

	m.Wg.Wait()
}
