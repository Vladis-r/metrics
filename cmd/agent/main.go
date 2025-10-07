package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/Vladis-r/metrics.git/cmd/agent/flags"
	"github.com/Vladis-r/metrics.git/internal/agent"
)

var (
	wg sync.WaitGroup
)

func main() {
	flag.Parse() // Parse command-line arguments

	fmt.Println("Start metrics agent...")

	gorutines := []func(*sync.WaitGroup){
		agent.GoUpdateMetrics,
		agent.GoReportMetics,
	}

	for _, gorutine := range gorutines {
		wg.Add(1)
		go gorutine(&wg)
	}

	fmt.Println("Press Ctrl+C to exit")
	fmt.Printf("Start with flags: -p %v, -r %v\n", flags.PollInterval, flags.ReportInterval)
	wg.Wait()
}
