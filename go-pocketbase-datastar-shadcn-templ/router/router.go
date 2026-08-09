package router

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"{{.gitserver}}/{{.owner}}/{{.name}}/assets"
	"{{.gitserver}}/{{.owner}}/{{.name}}/components"
)

func Setup(se *core.ServeEvent) error {
	// Serve assets from disk in development (so tailwind --watch output is
	// picked up without a restart) and from the embedded files in production.
	var staticFS fs.FS = assets.FS
	if os.Getenv("GO_ENV") != "production" {
		staticFS = os.DirFS("assets")
	}
	se.Router.GET("/assets/{path...}", apis.Static(staticFS, false))

	// shadcn-templ component script bundle. The rendered <script> src carries
	// a content hash (/components/shadcn-templ-<hash>.js), so the handler is
	// mounted on the whole prefix; it 404s any unknown name.
	se.Router.GET("/components/{path...}", func(e *core.RequestEvent) error {
		components.ScriptsHandler().ServeHTTP(e.Response, e.Request)
		return nil
	})

	se.Router.GET("/health", func(e *core.RequestEvent) error {
		return e.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	csrf := se.Router.Group("")
	csrf.Bind(requireSameOrigin())

	csrf.GET("/login", loginGetHandler)
	csrf.POST("/login", loginPostHandler)
	csrf.POST("/logout", logoutHandler)

	protected := csrf.Group("")
	protected.Bind(requireAuthWithRedirect())
	protected.GET("/", dashboardHandler)
	protected.GET("/settings", settingsGetHandler)
	protected.POST("/settings/password", settingsPasswordHandler)

	return se.Next()
}
