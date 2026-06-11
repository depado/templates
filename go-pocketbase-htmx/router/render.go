package router

import (
	"github.com/a-h/templ"
	"github.com/pocketbase/pocketbase/core"
)

func render(e *core.RequestEvent, component templ.Component) error {
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return component.Render(e.Request.Context(), e.Response)
}
