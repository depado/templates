package router

import (
	"github.com/a-h/templ"
	"github.com/pocketbase/pocketbase/core"
)

// render writes a templ component as an HTML response.
// Used both for full pages and for Datastar fragment responses: Datastar
// reads "text/html" and morphs the returned elements into the DOM by ID.
func render(e *core.RequestEvent, component templ.Component) error {
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	return component.Render(e.Request.Context(), e.Response)
}
