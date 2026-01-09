package server

import (
	"encoding/json"
	"os"
	"time"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"go.uber.org/zap"
)

// SaveMetricsToFile - save metrics in file.
func SaveMetricsToFile(s *models.MemStorage) {
	ticker := time.NewTicker(time.Duration(s.C.StoreInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		saveMetricsToFileLogic(s)
	}
}

// saveMetricsToFileLogic - logic for SaveMetricsToFile.
func saveMetricsToFileLogic(s *models.MemStorage) {
	s.Mu.RLock()
	metrics := s.Store
	s.Mu.RUnlock()

	listOfMetrics := make([]models.Metric, len(metrics))
	idx := 0
	for _, v := range metrics {
		listOfMetrics[idx] = v
		idx++
	}

	tmpFile := s.C.FileStoragePath + ".tmp"
	file, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		s.Log.Error("Error: Cant open file for save metrics!", zap.Error(err))
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(listOfMetrics)
	if err != nil {
		s.Log.Error("Error: cant write metrics in file!", zap.Error(err))
	}
	file.Close()

	if err := os.Rename(tmpFile, s.C.FileStoragePath); err != nil {
		s.Log.Error("Error: cant rename tmp file with metrics!", zap.Error(err))
	}
	s.Log.Info("Metrics succesfull save in file")
}

// LoadMetricsFromFile - load metrics from file and save in Store.
func LoadMetricsFromFile(s *models.MemStorage) {
	if !s.C.IsRestore {
		s.Log.Info("Skip load metrics from file", zap.String("Path", s.C.FileStoragePath))
		return
	}
	data, err := os.ReadFile(s.C.FileStoragePath)
	if err != nil {
		s.Log.Error("Error: cant read file with metrics!", zap.Error(err))
	}
	metrics := []models.Metric{}
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		s.Log.Error("Error: cant unmarshal json with metrics!", zap.Error(err))
	}

	// save metrics in Store.
	s.Mu.RLock()
	for i := range metrics {
		s.Store[metrics[i].ID] = metrics[i]
	}
	s.Mu.RUnlock()
	s.Log.Info("Metrics load from file", zap.String("Path", s.C.FileStoragePath))
}
