package server

import (
	"context"
	"encoding/json"
	"os"
	"time"

	models "github.com/Vladis-r/metrics.git/internal/model"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// SaveMetricToDb - save metrics to database every tick Conf.StoreInterval .
func SaveMetricToDb(s *models.MemStorage) {
	ticker := time.NewTicker(time.Duration(s.Conf.StoreInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		SaveMetricToDbLogic(s)
	}
}

// SaveMetricToDbLogic - logic for save metrics to database.
func SaveMetricToDbLogic(s *models.MemStorage) {
	s.Mu.RLock()
	metrics := s.Store
	s.Mu.RUnlock()

	if s.Conf.DatabaseDsn == "" {
		s.Log.Info("Database DSN not provided, skipping save to database")
		return
	}

	conn, err := pgx.Connect(context.Background(), s.Conf.DatabaseDsn)
	if err != nil {
		s.Log.Error("Failed to connect to database", zap.Error(err))
		return
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(context.Background())
	if err != nil {
		s.Log.Error("Failed to begin transaction", zap.Error(err))
		return
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Prepare(context.Background(), "upsert_metric",
		`INSERT INTO metrics (name, type, value, delta, hash) 
		VALUES ($1, $2, $3, $4, $5) 
		ON CONFLICT (name) 
		DO UPDATE SET value = EXCLUDED.value, delta = EXCLUDED.delta, hash = EXCLUDED.hash`,
	)
	if err != nil {
		s.Log.Error("Failed to prepare statement", zap.Error(err))
		return
	}
	for _, metric := range metrics {
		var value *float64
		var delta *int64
		switch metric.MType {
		case models.Gauge:
			value = metric.Value
		case models.Counter:
			delta = metric.Delta
		default:
			s.Log.Info("Unexpected type", zap.String("MType", metric.MType))
			continue
		}
		_, err = tx.Exec(context.Background(), "upsert_metric", metric.ID, metric.MType, value, delta, metric.Hash)
		if err != nil {
			s.Log.Error("Failed to save metric to database", zap.String("metric", metric.ID), zap.Error(err))
			return
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		s.Log.Error("Failed to commit transaction", zap.Error(err))
		return
	}

	s.Log.Info("Metrics successfully saved to database", zap.Int("count", len(metrics)))
}

// SaveMetricsToFile - save metrics in file.
func SaveMetricsToFile(s *models.MemStorage) {
	ticker := time.NewTicker(time.Duration(s.Conf.StoreInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		SaveMetricsToFileLogic(s)
	}
}

// SaveMetricsToFileLogic - logic for SaveMetricsToFile.
func SaveMetricsToFileLogic(s *models.MemStorage) {
	s.Mu.RLock()
	metrics := s.Store
	s.Mu.RUnlock()

	listOfMetrics := make([]models.Metric, len(metrics))
	idx := 0
	for _, v := range metrics {
		listOfMetrics[idx] = v
		idx++
	}

	tmpFile := s.Conf.FileStoragePath + ".tmp"
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

	if err := os.Rename(tmpFile, s.Conf.FileStoragePath); err != nil {
		s.Log.Error("Error: cant rename tmp file with metrics!", zap.Error(err))
	}
	s.Log.Info("Metrics succesfull save in file")
}

// LoadMetricsFromFile - load metrics from file and save in Store.
func LoadMetricsFromFile(s *models.MemStorage) {
	if !s.Conf.IsRestore {
		s.Log.Info("Skip load metrics from file", zap.String("Path", s.Conf.FileStoragePath))
		return
	}
	data, err := os.ReadFile(s.Conf.FileStoragePath)
	if err != nil {
		s.Log.Error("Error: cant read file with metrics!", zap.Error(err))
	}
	metrics := []models.Metric{}
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		s.Log.Error("Error: cant unmarshal json with metrics!", zap.Error(err))
	}

	s.Mu.Lock()
	for i := range metrics {
		s.Store[metrics[i].ID] = metrics[i]
	}
	s.Mu.Unlock()
	s.Log.Info("Metrics load from file", zap.String("Path", s.Conf.FileStoragePath))
}

// LoadMetricsFromDatabase - load metrics from database if DSN is provided
func LoadMetricsFromDatabase(s *models.MemStorage) {
	if s.Conf.DatabaseDsn == "" {
		s.Log.Info("Database DSN not provided, skipping load from database")
		return
	}

	conn, err := pgx.Connect(context.Background(), s.Conf.DatabaseDsn)
	if err != nil {
		s.Log.Error("Failed to connect to database", zap.Error(err))
		return
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT name, type, value, delta, hash FROM metrics")
	if err != nil {
		s.Log.Error("Failed to query metrics from database", zap.Error(err))
		return
	}
	defer rows.Close()

	metrics := make(map[string]models.Metric)
	for rows.Next() {
		var name, mType, hash string
		var value *float64
		var delta *int64

		err = rows.Scan(&name, &mType, &value, &delta, &hash)
		if err != nil {
			s.Log.Error("Failed to scan metric row", zap.Error(err))
			return
		}
		metric := models.Metric{
			ID:    name,
			MType: mType,
			Hash:  hash,
		}
		switch mType {
		case models.Gauge:
			metric.Value = value
		case models.Counter:
			metric.Delta = delta
		default:
			continue
		}
		metrics[name] = metric
	}
	if err = rows.Err(); err != nil {
		s.Log.Error("Error during iteration over metric rows", zap.Error(err))
		return
	}
	s.Mu.Lock()
	for name, metric := range metrics {
		s.Store[name] = metric
	}
	s.Mu.Unlock()
	s.Log.Info("Metrics successfully loaded from database", zap.Int("count", len(metrics)))
}
