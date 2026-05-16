---
if: gin
---
package server

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"

	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd"
)

// NewCors generates a new cors config
func NewCors(c *cmd.Conf, l *slog.Logger) (*cors.Config, error) {
	if !c.Server.Cors.Enabled {
		return nil, nil
	}
	cc := &cors.Config{
		AllowCredentials: true,
		MaxAge:           50 * time.Second,
		AllowMethods:     c.Server.Cors.Methods,
		AllowHeaders:     c.Server.Cors.Headers,
		ExposeHeaders:    c.Server.Cors.Expose,
	}

	switch {
	case len(c.Server.Cors.Origins) > 0:
		cc.AllowOrigins = c.Server.Cors.Origins
	case c.Server.Cors.All:
		cc.AllowAllOrigins = true
	default:
		return nil, fmt.Errorf("all origins disabled but no allowed origin provided")
	}

	l.Debug("CORS configuration",
		"methods", cc.AllowMethods,
		"headers", cc.AllowHeaders,
		"origins", cc.AllowOrigins,
		"all_origins", cc.AllowAllOrigins,
	)

	return cc, nil
}
