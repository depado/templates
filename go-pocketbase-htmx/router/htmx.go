package router

import (
	"github.com/donseba/go-htmx"
	"github.com/pocketbase/pocketbase/core"
)

var htmxService = htmx.New()

func NewHandler(e *core.RequestEvent) *htmx.Handler {
	return htmxService.NewHandler(e.Response, e.Request)
}
