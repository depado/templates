package router

import (
	"github.com/pocketbase/pocketbase/core"
)

func LoggerMiddleware(e *core.RequestEvent) error {
	e.App.Logger().Info("request",
		"method", e.Request.Method,
		"path", e.Request.URL.Path,
		"ip", e.RealIP(),
	)
	return e.Next()
}

func RecoverMiddleware(e *core.RequestEvent) error {
	defer func() {
		if r := recover(); r != nil {
			e.App.Logger().Error("panic recovered", "error", r)
			e.InternalServerError("Internal server error", nil)
		}
	}()
	return e.Next()
}
