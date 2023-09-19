---
if: gin
---
package router

import (
	"fmt"
	"net/http"
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd"
)

// Router is a simple struct grouping the logger, configuration and gin engine.
type Router struct {
	conf *cmd.Conf
	l    *zerolog.Logger
	e    *gin.Engine
}

// New will create a new Router.
func New(c *cmd.Conf, l *zerolog.Logger, e *gin.Engine) Router {
	return Router{c, l, e}
}

// SetRoutes will add the necessary routes to the router.
func (r Router) setRoutes() {
	r.e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// Listen will start the server and will listen to incoming requests until asked
// to stop, in which case it will gracefully shutdown with a 5 seconds timeout.
func (r Router) Listen() {
	// Set routes
	r.setRoutes()

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", r.conf.Server.Host, r.conf.Server.Port),
		Handler: r.e,
	}

	// Initializing the server in a goroutine so that it won't block the
	// graceful shutdown handling below
	go func() {
		r.l.Info().Str("host", r.conf.Server.Host).Int("port", r.conf.Server.Port).Bool("cors", r.conf.Server.Cors.Enabled).Msg("listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.l.Fatal().Err(err).Msg("listen")
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	r.l.Info().Msg("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		r.l.Fatal().Err(err).Msg("server forced to shutdown")
	}
	r.l.Info().Msg("server exiting")
}
