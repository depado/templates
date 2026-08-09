package router

import (
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"

	"{{.gitserver}}/{{.owner}}/{{.name}}/views/pages"
)

const (
	authCookieName   = "pb_auth"
	authMiddlewareID = "requireAuthWithRedirect"
)

func requireAuthWithRedirect() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id: authMiddlewareID,
		Func: func(e *core.RequestEvent) error {
			if e.Auth == nil {
				cookie, err := e.Request.Cookie(authCookieName)
				if err == nil && cookie.Value != "" {
					e.Auth, _ = e.App.FindAuthRecordByToken(cookie.Value, core.TokenTypeAuth)
				}
			}
			if e.Auth == nil {
				if isDatastar(e) {
					// A 302 would morph the login page into the current fragment;
					// stream the Datastar redirect instead (patches a script that
					// sets window.location).
					return newSSE(e).Redirect("/login")
				}
				return e.Redirect(http.StatusFound, "/login")
			}
			return e.Next()
		},
	}
}

func setAuthCookie(e *core.RequestEvent, record *core.Record) error {
	token, err := record.NewAuthToken()
	if err != nil {
		return err
	}
	http.SetCookie(e.Response, &http.Cookie{
		Name:     authCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   e.IsTLS(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((14 * 24 * time.Hour).Seconds()),
	})
	return nil
}

func clearAuthCookie(e *core.RequestEvent) {
	http.SetCookie(e.Response, &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func loginGetHandler(e *core.RequestEvent) error {
	cookie, err := e.Request.Cookie(authCookieName)
	if err == nil && cookie.Value != "" {
		if _, err := e.App.FindAuthRecordByToken(cookie.Value, core.TokenTypeAuth); err == nil {
			return e.Redirect(http.StatusFound, "/")
		}
	}
	return render(e, pages.LoginPage())
}

func loginPostHandler(e *core.RequestEvent) error {
	email := e.Request.PostFormValue("email")
	password := e.Request.PostFormValue("password")

	record, err := e.App.FindAuthRecordByEmail("users", email)
	if err != nil || !record.ValidatePassword(password) {
		// Patch the login form's message element in place; typed values are
		// preserved by the morph.
		return newSSE(e).PatchElementTempl(pages.LoginMessage("Invalid email or password."))
	}

	if err := setAuthCookie(e, record); err != nil {
		return e.InternalServerError("Failed to set auth cookie", err)
	}

	return newSSE(e).Redirect("/")
}

func logoutHandler(e *core.RequestEvent) error {
	clearAuthCookie(e)
	return newSSE(e).Redirect("/login")
}
