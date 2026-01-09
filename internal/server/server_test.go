package server

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/Vladis-r/metrics.git/cmd/config"
	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/Vladis-r/metrics.git/internal/utils"
	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

func newTestLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func createTestMemStorageSaveFile(filePath string, metrics map[string]models.Metric) *models.MemStorage {
	return &models.MemStorage{
		C: &config.ConfigServer{
			FileStoragePath: filePath,
		},
		Store: metrics,
		Mu:    sync.RWMutex{},
		Log:   newTestLogger(),
	}
}

func TestSaveMetricsToFile(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		store    map[string]models.Metric
	}{
		{
			name:     "Test 1. Successful save metrics to file",
			filePath: "test_metrics.json",
			store: map[string]models.Metric{
				"gauge1": {
					ID:    "gauge1",
					MType: "gauge",
					Value: utils.Float64Ptr(3.14),
				},
				"counter1": {
					ID:    "counter1",
					MType: "counter",
					Delta: utils.Int64Ptr(42),
				},
			},
		},
		{
			name:     "Test 2. Successful save empty metrics to file",
			filePath: "test_metrics.json",
			store:    map[string]models.Metric{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "metrics.json")

			storage := createTestMemStorageSaveFile(outputFile, tt.store)
			saveMetricsToFileLogic(storage)

			require.FileExists(t, outputFile, "Output file should exist")

			data, err := os.ReadFile(outputFile)
			require.NoError(t, err, "Should be able to read the output file")

			var savedMetrics []models.Metric
			err = json.Unmarshal(data, &savedMetrics)
			require.NoError(t, err, "Should be valid JSON")

			require.Len(t, savedMetrics, len(tt.store), "Should save 2 metrics")

			metricsMap := make(map[string]models.Metric)
			for _, m := range savedMetrics {
				metricsMap[m.ID] = m
			}
			if len(tt.store) > 0 {
				require.Equal(t, 3.14, *metricsMap["gauge1"].Value)
				require.Equal(t, "gauge", metricsMap["gauge1"].MType)
				require.Equal(t, int64(42), *metricsMap["counter1"].Delta)
				require.Equal(t, "counter", metricsMap["counter1"].MType)
			}

			tmpFile := outputFile + ".tmp"
			require.NoFileExists(t, tmpFile, "Temporary file should be removed")
		})
	}

	errorTests := []struct {
		name     string
		filePath string
		store    map[string]models.Metric
	}{
		{
			name:     "error while save metrics to file",
			filePath: "test_metrics.json",
			store:    map[string]models.Metric{},
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := "invalid/path"
			outputFile := filepath.Join(tmpDir, "metrics.json")

			storage := createTestMemStorageSaveFile(outputFile, tt.store)
			saveMetricsToFileLogic(storage)
			require.NoFileExists(t, outputFile, "File exist with invalid path")
		})
	}
}

func createTestMemStorageLoadFile(filePath, isRestore string, metrics map[string]models.Metric) *models.MemStorage {
	return &models.MemStorage{
		C: &config.ConfigServer{
			FileStoragePath: filePath,
			IsRestore:       isRestore == "true",
		},
		Store: metrics,
		Mu:    sync.RWMutex{},
		Log:   newTestLogger(),
	}
}

func TestLoadMetricsFromFile(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		isRestore string
		store     []models.Metric
	}{
		{
			name:      "Test 1. Successful load metrics from file",
			filePath:  "test_metrics.json",
			isRestore: "true",
			store: []models.Metric{
				{
					ID:    "gauge1",
					MType: "gauge",
					Value: utils.Float64Ptr(3.14),
				},
				{
					ID:    "counter1",
					MType: "counter",
					Delta: utils.Int64Ptr(42),
				},
			},
		},
		{
			name:      "Test 2. Skip load metrics from files - isRestore == false.",
			filePath:  "test_metrics.json",
			isRestore: "false",
			store:     []models.Metric{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filepath := filepath.Join(tmpDir, "metrics.json")

			data, err := json.Marshal(tt.store)
			require.NoError(t, err)
			err = os.WriteFile(filepath, data, 0600)
			require.NoError(t, err)

			storage := createTestMemStorageLoadFile(filepath, tt.isRestore, map[string]models.Metric{})
			LoadMetricsFromFile(storage)

			if !storage.C.IsRestore {
				require.Empty(t, storage.Store, "Should not load metrics when IsRestore is false")
			}
			if len(tt.store) > 0 {
				require.Equal(t, 3.14, *storage.Store["gauge1"].Value)
				require.Equal(t, "gauge", storage.Store["gauge1"].MType)
				require.Equal(t, int64(42), *storage.Store["counter1"].Delta)
				require.Equal(t, "counter", storage.Store["counter1"].MType)
			}
		})
	}
}
