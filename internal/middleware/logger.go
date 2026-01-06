package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// responseWriter — обёртка над gin.ResponseWriter для захвата тела ответа.
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write — перехватываем запись тела ответа
func (w *responseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// Loger - wraps HTTP handlers for logger.
func Logger(l *zap.Logger) gin.HandlerFunc {
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
		if c.Request.URL.Path == "/" {
			body = []byte{}
		}

		writer := &responseWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = writer

		l.Info("Request",
			zap.String("url", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("body", string(body)),
		)

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

		bodyResp := ""
		if c.Request.URL.Path != "/" {
			bodyResp = writer.body.String()
		}

		l.Info("Response",
			zap.Int("status", statusCode),
			zap.String("body", bodyResp),
			zap.Duration("duration", duration),
			zap.Int64("size", int64(size)))
	}
}

// InitLogger - инициализация логгера.
func InitLogger() (logger *zap.Logger, err error) {
	logger, err = zap.NewProduction()
	return logger, err
}
