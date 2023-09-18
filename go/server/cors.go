---
if: gin
---
package server

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/rs/zerolog"

	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd"
)

// NewCorsConfig generates a new cors config
func NewCors(c *cmd.Conf, l zerolog.Logger) (*cors.Config, error) {
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
		return nil, fmt.Errorf("all all origins disabled but no allowed origin provided")
	}

	l.Debug().
		Strs("methods", cc.AllowMethods).
		Strs("headers", cc.AllowHeaders).
		Strs("headers", cc.AllowHeaders).
		Strs("origins", cc.AllowOrigins).
		Bool("all_origins", cc.AllowAllOrigins).
		Msg("CORS configuration")

	return cc, nil
}
