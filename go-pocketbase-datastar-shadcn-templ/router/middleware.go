package router

import (
	"net/http"
	"net/url"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
)

// requireSameOrigin is a lightweight CSRF protection for cookie-authenticated
// state-changing endpoints: for non-safe methods it checks that the Origin
// header, when present, matches the request host. Browsers always send Origin
// on cross-origin POSTs and on same-origin fetches, so a mismatch means a
// cross-site request. Requests without an Origin (curl, server-to-server,
// mobile SDKs) are allowed: the pb_auth cookie is HttpOnly + SameSite=Lax,
// which cross-site browser requests cannot attach anyway.
func requireSameOrigin() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id: "requireSameOrigin",
		Func: func(e *core.RequestEvent) error {
			switch e.Request.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions:
				return e.Next()
			}

			origin := e.Request.Header.Get("Origin")
			if origin == "" {
				return e.Next()
			}

			u, err := url.Parse(origin)
			if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host != e.Request.Host {
				return e.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
			}

			return e.Next()
		},
	}
}
