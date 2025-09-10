package main

import (
	"fmt"
	"sync"

	"github.com/Vladis-r/metrics.git/internal/agent"
)

func main() {
	fmt.Println("Start metrics agent...")
	var (
		wg sync.WaitGroup
	)
	gorutines := []func(*sync.WaitGroup){
		agent.GoUpdateMetrics,
		agent.GoReportMetics,
	}
	for _, gorutine := range gorutines {
		wg.Add(1)
		go gorutine(&wg)
	}

	fmt.Println("Press Ctrl+C to exit")
	wg.Wait()
}
