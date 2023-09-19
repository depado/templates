---
if: gin
---
package server

import (
	"github.com/Depado/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd"
)

// setMode is used to set the proper gin mode.
func setMode(mode string, l *zerolog.Logger) {
	switch mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		l.Warn().Str("mode", mode).Msg("unknown gin mode, fallback to release")
		gin.SetMode(gin.ReleaseMode)
	}
}

// NewGinEngine will configure and return a new gin engine.
func NewGinEngine(c *cmd.Conf, l *zerolog.Logger, cc *cors.Config) *gin.Engine {
	setMode(c.Server.Mode, l)
	r := gin.New()

	// Setup instrumentation if configured
	if c.Server.Instrument {
		p := ginprom.New(ginprom.Engine(r))
		r.Use(p.Instrument())
	}

	// Setup logging
	if c.Server.UnifiedLogger {
		r.Use(ZerologLogger(l))
	} else {
		r.Use(gin.Logger())
	}

	// Recovers on panic
	r.Use(gin.Recovery())

	// Setup cors is configured
	if cc != nil {
		r.Use(cors.New(*cc))
	}

	return r
}
