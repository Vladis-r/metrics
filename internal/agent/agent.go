package agent

import (
	"crypto/rand"
	"fmt"
	"maps"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	PollInterval       = 2 * time.Second  // time for update metrics
	ReportInterval     = 10 * time.Second // time for report metrics to server
	port               = ":8080"
	server             = "http://localhost"
	routeUpdateMetrics = "/update"
)

// MetricsRuntimeMap - map with runTime metrics. "name-type": "value"
var MetricsMap = map[string]string{
	// Runtime
	"Alloc-counter":        "value",
	"BuckHashSys-counter":  "value",
	"Frees-counter":        "value",
	"GCCPUFraction-gauge":  "value",
	"GCSys-counter":        "value",
	"HeapAlloc-counter":    "value",
	"HeapIdle-counter":     "value",
	"HeapInuse-counter":    "value",
	"HeapObjects-counter":  "value",
	"HeapReleased-counter": "value",
	"HeapSys-counter":      "value",
	"LastGC-counter":       "value",
	"Lookups-counter":      "value",
	"MCacheInuse-counter":  "value",
	"MCacheSys-counter":    "value",
	"MSpanInuse-counter":   "value",
	"MSpanSys-counter":     "value",
	"Mallocs-counter":      "value",
	"NextGC-counter":       "value",
	"NumForcedGC-counter":  "value",
	"NumGC-counter":        "value",
	"OtherSys-counter":     "value",
	"PauseTotalNs-counter": "value",
	"StackInuse-counter":   "value",
	"StackSys-counter":     "value",
	"Sys-counter":          "value",
	"TotalAlloc-counter":   "value",
	// Custom
	"PollCount-counter": "value",
	"RandomValue-gauge": "value",
}

var mu sync.Mutex

// GoReportMetics - func
func GoReportMetics(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(ReportInterval)
	defer ticker.Stop()

	for t := range ticker.C {
		mu.Lock()
		newMap := make(map[string]string, len(MetricsMap))
		maps.Copy(newMap, MetricsMap)
		mu.Unlock()
		go sendMetrics(newMap)
		fmt.Println("Send metrics. Tick at ", t)
	}
}

func GoUpdateMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	for t := range ticker.C {
		mu.Lock()
		updateMetrics()
		mu.Unlock()
		fmt.Println("Update metrics. Tick at ", t)
	}
}

func updateMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	for key := range MetricsMap {
		switch key {
		case "PollCount-counter":
			MetricsMap[key] = getPollCountertMetric(MetricsMap[key])
		case "RandomValue-gauge":
			MetricsMap[key] = getRandomValueMetric()
		default:
			MetricsMap[key] = getRunTimeMetrics(key, memStats)
		}
	}
}

func sendMetrics(metricsMap map[string]string) {
	for key, metricValue := range metricsMap {
		splittedString := strings.Split(key, "-")
		metricName, metricType := splittedString[0], splittedString[1]
		fullUrl := fmt.Sprintf("%s%s%s/%s/%s/%s", server, port, routeUpdateMetrics, metricType, metricName, metricValue)
		resp, err := http.Post(fullUrl, "text/plain", nil)
		if err != nil {
			e := fmt.Errorf("error send metrics: %w", err)
			fmt.Println(e)
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Println(resp.StatusCode)
			fmt.Println(fullUrl)
		}
	}
}

func getRandomValueMetric() string {
	randInt, _ := rand.Prime(rand.Reader, 64)
	randFloat, _ := randInt.Float64()
	return fmt.Sprintf("%f", randFloat)
}

func getPollCountertMetric(counter string) string {
	v, err := strconv.Atoi(counter)
	if err != nil {
		v = -1
	}
	return fmt.Sprintf("%d", v+1)
}

func getRunTimeMetrics(key string, memStats runtime.MemStats) string {
	name := strings.Split(key, "-")[0]
	value := fmt.Sprintf("%v", reflect.ValueOf(memStats).FieldByName(name))
	return value
}
