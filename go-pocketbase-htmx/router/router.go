package router

import (
	"net/http"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"{{.gitserver}}/{{.owner}}/{{.name}}/assets"
)

func Setup(se *core.ServeEvent) error {
	se.Router.GET("/assets/{path...}", apis.Static(assets.FS, false))

	se.Router.GET("/health", func(e *core.RequestEvent) error {
		return e.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	se.Router.GET("/login", loginGetHandler)
	se.Router.POST("/login", loginPostHandler)
	se.Router.GET("/logout", logoutHandler)

	protected := se.Router.Group("")
	protected.Bind(requireAuthWithRedirect())
	protected.GET("/", dashboardHandler)
	protected.GET("/settings", settingsGetHandler)
	protected.POST("/settings/password", settingsPasswordHandler)

	return se.Next()
}
