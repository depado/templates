package router

import (
	"github.com/pocketbase/pocketbase/core"
)

// LoggerMiddleware logs incoming requests.
func LoggerMiddleware(e *core.RequestEvent) error {
	e.App.Logger().Info("request",
		"method", e.Request.Method,
		"path", e.Request.URL.Path,
		"ip", e.RealIP(),
	)
	return e.Next()
}

// RecoverMiddleware recovers from panics and logs the error.
func RecoverMiddleware(e *core.RequestEvent) error {
	defer func() {
		if r := recover(); r != nil {
			e.App.Logger().Error("panic recovered", "error", r)
			e.InternalServerError("Internal server error", nil)
		}
	}()
	return e.Next()
}
