package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	w.Header().Set("Content-Encoding", "gzip")
	return w.writer.Write(data)
}

func (w *gzipWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

func (w *gzipWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			defer gz.Close()
			c.Request.Body = io.NopCloser(gz)
		}
		if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			gz := gzip.NewWriter(c.Writer)
			defer gz.Close()

			c.Writer = &gzipWriter{c.Writer, gz}

			c.Header("Vary", "Accept-Encoding")
		}
		c.Next()
		if gw, ok := c.Writer.(*gzipWriter); ok {
			gw.Flush()
		}
	}
}
