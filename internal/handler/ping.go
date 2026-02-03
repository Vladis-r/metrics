package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

func Ping(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		err := db.PingContext(ctx)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "detail": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	}
}
