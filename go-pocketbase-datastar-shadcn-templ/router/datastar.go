package router

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/starfederation/datastar-go/datastar"
)

// isDatastar reports whether the request was made by a Datastar backend action.
// Datastar sends the "Datastar-Request: true" header on all @get/@post/... fetches.
func isDatastar(e *core.RequestEvent) bool {
	return e.Request.Header.Get("Datastar-Request") == "true"
}

// newSSE creates a ServerSentEventGenerator writing to the current request.
func newSSE(e *core.RequestEvent) *datastar.ServerSentEventGenerator {
	return datastar.NewSSE(e.Response, e.Request)
}
