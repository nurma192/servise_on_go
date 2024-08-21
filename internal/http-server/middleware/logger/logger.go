package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

func New(log *slog.Logger) gin.HandlerFunc {
	log = log.With(
		slog.String("component", "middleware/logger"),
	)
	log.Info("logger middleware enabled")

	return func(c *gin.Context) {
		// Generate a unique request ID
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)

		// Create a logger entry with request details
		entry := log.With(
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("remote_addr", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("request_id", requestID),
		)

		// Record the start time
		startTime := time.Now()

		// Process request
		c.Next()

		// After request is processed
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		dataLength := c.Writer.Size()

		entry.Info("request completed",
			slog.Int("status", statusCode),
			slog.Int("bytes", dataLength),
			slog.String("duration", latency.String()),
		)
	}
}
