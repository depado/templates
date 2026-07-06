---
if: gin
---
package server

import (
	"context"
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

// multiHandler fans out log records to multiple slog.Handler backends,
// e.g. stderr (tint) and OTLP. Each handler receives its own clone of the
// record.
type multiHandler struct {
	handlers []slog.Handler
}

func (h multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return multiHandler{handlers: handlers}
}

func (h multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return multiHandler{handlers: handlers}
}
