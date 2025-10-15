package main

import (
	"fmt"

	"github.com/Vladis-r/metrics.git/cmd/config"
	"github.com/Vladis-r/metrics.git/internal/agent"
	models "github.com/Vladis-r/metrics.git/internal/model"
)

func main() {
	c := config.GetConfig()     // Parse command-line arguments.
	m := models.NewMetricsMap() // Init client and map for metrics.

	fmt.Println("Start metrics agent...")
	fmt.Printf("With config:\n PollInterval: %v\n ReportInterval: %v\n\n", c.PollInterval, c.ReportInterval)

	goroutines := []func(*models.MetricsMap, *config.Config){
		agent.GoUpdateMetrics,
		agent.GoReportMetrics,
	}

	for _, goroutine := range goroutines {
		m.Wg.Add(1)
		go goroutine(m, c)
	}

	fmt.Println("Press Ctrl+C to exit")

	m.Wg.Wait()
}
