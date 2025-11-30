package logger

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Log *zap.Logger

// Middleware - wraps HTTP handlers for logger.
func Middleware(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var body []byte
		if c.Request.Body != nil {
			var err error
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				l.Error("Failed to read request body", zap.Error(err))
				c.Next()
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		c.Next()

		statusCode := c.Writer.Status()
		if statusCode == 0 {
			statusCode = 200
		}
		size := c.Writer.Size()
		if size < 0 {
			size = 0
		}

		duration := time.Since(start)

		l.Info("Request",
			zap.String("url", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Duration("duration", duration),
			zap.String("body", truncateString(string(body), 1024)),
		)
		l.Info("Response",
			zap.Int("status", statusCode),
			zap.Int64("size", int64(size)))
	}
}

// InitLogger - инициализация логгера.
func InitLogger() (err error) {
	Log, err = zap.NewProduction()
	if err != nil {
		Log.Panic("Cant init logger!")
	}

	return nil
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "... [truncated]"
}
