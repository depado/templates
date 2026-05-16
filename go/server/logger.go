---
if: gin
---
package server

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// SlogLogger logs a gin HTTP request with the provided slog logger.
func SlogLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		latency := time.Since(start)
		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}

		if raw != "" {
			path = path + "?" + raw
		}

		status := c.Writer.Status()
		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		attrs := []any{
			"client_ip", c.ClientIP(),
			"method", c.Request.Method,
			"status", status,
			"body_size", c.Writer.Size(),
			"path", path,
			"latency", latency.String(),
		}

		if errMsg != "" {
			attrs = append(attrs, "error", errMsg)
		}

		if status >= 500 {
			logger.Error("request", attrs...)
		} else {
			logger.Info("request", attrs...)
		}
	}
}
