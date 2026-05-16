package router

import (
	"net/http"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// Setup configures all routes and middlewares for the application.
func Setup(se *core.ServeEvent) error {
	// Public routes
	se.Router.GET("/health", healthHandler)

	// API routes group
	api := se.Router.Group("/api/v1")
	api.GET("/hello/{name}", helloHandler)

	// Protected routes example
	protected := se.Router.Group("/api/v1/protected")
	protected.Bind(apis.RequireAuth())
	protected.GET("/me", meHandler)

	return se.Next()
}

// healthHandler returns the health status of the application.
func healthHandler(e *core.RequestEvent) error {
	return e.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// helloHandler is a simple example handler.
func helloHandler(e *core.RequestEvent) error {
	name := e.Request.PathValue("name")
	return e.JSON(http.StatusOK, map[string]string{"message": "Hello " + name})
}

// meHandler returns the currently authenticated user.
func meHandler(e *core.RequestEvent) error {
	record := e.Auth
	if record == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}
	return e.JSON(http.StatusOK, record)
}
