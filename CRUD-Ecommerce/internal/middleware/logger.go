package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		requestID := uuid.New().String()

		c.Set("request_id", requestID)

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Log request details
		c.Writer.Header().Set("X-Request-ID", requestID)

		// You can add more sophisticated logging here (e.g., to a file or external service)
		if statusCode >= 500 {
			// Log errors
			c.Error(nil)
		}

		// Simple console log
		println("["+method+"]", path, statusCode, duration.String())
	}
}
