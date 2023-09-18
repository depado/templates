---
if: gin
---
package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd"
)

type Router struct {
	conf *cmd.Conf
	l    zerolog.Logger
	e    *gin.Engine
}

func New(c *cmd.Conf, l zerolog.Logger, e *gin.Engine) Router {
	return Router{c, l, e}
}

func (r Router) SetRoutes() {
	r.e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func (r Router) Listen() error {
	r.SetRoutes()
	r.l.Info().Str("host", r.conf.Server.Host).Int("port", r.conf.Server.Port).Bool("cors", r.conf.Server.Cors.Enabled).Msg("listening")
	if err := r.e.Run(fmt.Sprintf("%s:%d", r.conf.Server.Host, r.conf.Server.Port)); err != nil {
		return fmt.Errorf("running router: %w", err)
	}
	return nil
}
