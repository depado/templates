package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

const (
	testEmail    = "test@example.com"
	testPassword = "1234567890"
)

// setupRouter builds a test PocketBase app with the router wired the same way
// as main.go does at runtime.
func setupRouter(t *testing.T) (*tests.TestApp, http.Handler) {
	t.Helper()

	app, err := tests.NewTestApp()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { app.Cleanup() })

	baseRouter, err := apis.NewRouter(app)
	if err != nil {
		t.Fatal(err)
	}

	serveEvent := new(core.ServeEvent)
	serveEvent.App = app
	serveEvent.Router = baseRouter

	if err := app.OnServe().Trigger(serveEvent, func(e *core.ServeEvent) error {
		return Setup(e)
	}); err != nil {
		t.Fatal(err)
	}

	mux, err := baseRouter.BuildMux()
	if err != nil {
		t.Fatal(err)
	}

	return app, mux
}

func doRequest(mux http.Handler, method, path string, opts ...func(*http.Request)) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, o := range opts {
		o(req)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func withDatastar() func(*http.Request) {
	return func(r *http.Request) {
		r.Header.Set("Datastar-Request", "true")
	}
}

func withOrigin(origin string) func(*http.Request) {
	return func(r *http.Request) {
		r.Header.Set("Origin", origin)
	}
}

func withForm(v url.Values) func(*http.Request) {
	return func(r *http.Request) {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		body := strings.NewReader(v.Encode())
		r.Body = io.NopCloser(body)
		r.ContentLength = int64(body.Len())
	}
}

// authCookie returns the pb_auth cookie for the seed user.
func authCookie(t *testing.T, app *tests.TestApp) *http.Cookie {
	t.Helper()
	record, err := app.FindAuthRecordByEmail("users", testEmail)
	if err != nil {
		t.Fatal(err)
	}
	token, err := record.NewAuthToken()
	if err != nil {
		t.Fatal(err)
	}
	return &http.Cookie{Name: authCookieName, Value: token}
}

func TestHealth(t *testing.T) {
	_, mux := setupRouter(t)

	rec := doRequest(mux, http.MethodGet, "/health")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestLoginPage(t *testing.T) {
	t.Setenv("GO_ENV", "production")
	_, mux := setupRouter(t)

	rec := doRequest(mux, http.MethodGet, "/login")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		`id="login-form"`,
		`data-on:submit="@post('/login', {contentType: 'form'})"`,
		`/assets/js/datastar.js`,
		`/components/shadcn-templ-`,
		`id="login-message"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("login page missing %q", want)
		}
	}
}

func TestStaticAssets(t *testing.T) {
	// Embedded files: tests run from the package dir, so the dev mode
	// (os.DirFS("assets")) cannot resolve the assets here.
	t.Setenv("GO_ENV", "production")
	_, mux := setupRouter(t)

	for _, path := range []string{"/assets/css/output.css", "/assets/js/datastar.js"} {
		rec := doRequest(mux, http.MethodGet, path)
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s: expected 200, got %d", path, rec.Code)
		}
	}
}

// TestScriptsBundle verifies the shadcn-templ component bundle is served both
// on the plain and the content-hashed URL the <script> tag renders. GO_ENV is
// pinned to production so the bundle comes from the module's embedded files
// (the dev mode reads the local components/ dir from the working directory,
// which tests cannot rely on).
func TestScriptsBundle(t *testing.T) {
	t.Setenv("GO_ENV", "production")
	_, mux := setupRouter(t)

	// The rendered page must reference a fetchable hashed bundle URL.
	page := doRequest(mux, http.MethodGet, "/login")
	src := regexp.MustCompile(`src="(/components/shadcn-templ-[^"]+)"`).FindStringSubmatch(page.Body.String())
	if len(src) != 2 {
		t.Fatalf("could not find hashed bundle src in %s", page.Body.String())
	}

	for _, path := range []string{src[1], "/components/shadcn-templ.js"} {
		rec := doRequest(mux, http.MethodGet, path)
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s: expected 200, got %d", path, rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/javascript") {
			t.Errorf("GET %s: unexpected content type %q", path, ct)
		}
	}

	rec := doRequest(mux, http.MethodGet, "/components/shadcn-templ.js")
	body := rec.Body.String()
	for _, want := range []string{"window.tui", "toast", "sidebar"} {
		if !strings.Contains(body, want) {
			t.Errorf("bundle missing %q", want)
		}
	}

	rec = doRequest(mux, http.MethodGet, "/components/shadcn-templ-unknown.js")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown bundle name, got %d", rec.Code)
	}
}

func TestProtectedRouteRedirectsToLogin(t *testing.T) {
	_, mux := setupRouter(t)

	rec := doRequest(mux, http.MethodGet, "/")
	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/login" {
		t.Fatalf("expected redirect to /login, got %q", loc)
	}
}

func TestProtectedRouteSSERedirectForDatastar(t *testing.T) {
	_, mux := setupRouter(t)

	rec := doRequest(mux, http.MethodGet, "/", withDatastar())
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 SSE, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Fatalf("expected text/event-stream, got %q", ct)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `window.location.href = "/login"`) {
		t.Fatalf("expected SSE redirect to /login, got: %s", body)
	}
}

func TestDashboardFullPage(t *testing.T) {
	app, mux := setupRouter(t)

	rec := doRequest(mux, http.MethodGet, "/", func(r *http.Request) {
		r.AddCookie(authCookie(t, app))
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{`id="main-content"`, `id="main-nav"`, `Welcome`} {
		if !strings.Contains(body, want) {
			t.Errorf("dashboard missing %q", want)
		}
	}
}

func TestFragmentNavigation(t *testing.T) {
	app, mux := setupRouter(t)

	rec := doRequest(mux, http.MethodGet, "/settings", withDatastar(), func(r *http.Request) {
		r.AddCookie(authCookie(t, app))
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		`<title id="page-title">`,
		`id="page-heading"`,
		`id="main-content"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("fragment missing %q", want)
		}
	}
	// The sidebar nav is rendered once on the full page; fragments must not
	// re-render it (avoids the tooltip/morph flash).
	if strings.Contains(body, `id="main-nav"`) {
		t.Error("fragment should not re-render the sidebar nav")
	}
	if strings.Contains(body, `<!DOCTYPE html>`) {
		t.Error("fragment response should not contain a full page")
	}
}

func TestLoginFlow(t *testing.T) {
	_, mux := setupRouter(t)

	// Wrong credentials: SSE patch with the error message.
	rec := doRequest(mux, http.MethodPost, "/login",
		withOrigin("http://example.com"),
		withForm(url.Values{"email": {testEmail}, "password": {"wrong"}}))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, "Invalid email or password.") {
		t.Fatalf("expected login error patch, got: %s", body)
	}

	// Correct credentials: SSE redirect + auth cookie.
	rec = doRequest(mux, http.MethodPost, "/login",
		withOrigin("http://example.com"),
		withForm(url.Values{"email": {testEmail}, "password": {testPassword}}))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `window.location.href = "/"`) {
		t.Fatalf("expected SSE redirect to /, got: %s", body)
	}
	var auth *http.Cookie
	for _, c := range rec.Result().Cookies() {
		if c.Name == authCookieName {
			auth = c
		}
	}
	if auth == nil || auth.Value == "" {
		t.Fatal("expected pb_auth cookie to be set")
	}
	if !auth.HttpOnly {
		t.Error("pb_auth cookie must be HttpOnly")
	}
	if auth.SameSite != http.SameSiteLaxMode {
		t.Errorf("expected SameSite=Lax, got %v", auth.SameSite)
	}

	// The cookie now grants access to protected routes.
	rec = doRequest(mux, http.MethodGet, "/", func(r *http.Request) {
		r.AddCookie(auth)
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 with auth cookie, got %d", rec.Code)
	}
}

func TestCSRFOriginMismatch(t *testing.T) {
	_, mux := setupRouter(t)

	rec := doRequest(mux, http.MethodPost, "/login",
		withOrigin("https://evil.example.com"),
		withForm(url.Values{"email": {testEmail}, "password": {testPassword}}))
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for cross-origin POST, got %d", rec.Code)
	}

	// Same-origin POST is processed (SSE response, not 403).
	rec = doRequest(mux, http.MethodPost, "/login",
		withOrigin("http://example.com"),
		withForm(url.Values{"email": {testEmail}, "password": {testPassword}}))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for same-origin POST, got %d", rec.Code)
	}

	// Missing Origin (curl/server-to-server) is allowed.
	rec = doRequest(mux, http.MethodPost, "/login",
		withForm(url.Values{"email": {testEmail}, "password": {testPassword}}))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 without Origin, got %d", rec.Code)
	}
}

func TestLogout(t *testing.T) {
	app, mux := setupRouter(t)

	rec := doRequest(mux, http.MethodPost, "/logout",
		withOrigin("http://example.com"),
		func(r *http.Request) {
			r.AddCookie(authCookie(t, app))
		})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, `window.location.href = "/login"`) {
		t.Fatalf("expected SSE redirect to /login, got: %s", body)
	}
	for _, c := range rec.Result().Cookies() {
		if c.Name == authCookieName && c.MaxAge >= 0 {
			t.Error("expected the auth cookie to be expired")
		}
	}
}

func TestPasswordChangeFlow(t *testing.T) {
	app, mux := setupRouter(t)

	post := func(current, newpw, confirm string) *httptest.ResponseRecorder {
		return doRequest(mux, http.MethodPost, "/settings/password",
			withOrigin("http://example.com"),
			func(r *http.Request) {
				r.AddCookie(authCookie(t, app))
			},
			withForm(url.Values{
				"current_password": {current},
				"new_password":     {newpw},
				"confirm_password": {confirm},
			}))
	}

	// Wrong current password.
	rec := post("wrong", "newpass123", "newpass123")
	if body := rec.Body.String(); !strings.Contains(body, "Current password is incorrect") {
		t.Fatalf("expected error message, got: %s", body)
	}

	// Mismatched confirmation.
	rec = post(testPassword, "newpass123", "otherpass")
	if body := rec.Body.String(); !strings.Contains(body, "Passwords do not match") {
		t.Fatalf("expected error message, got: %s", body)
	}

	// Too short.
	rec = post(testPassword, "short", "short")
	if body := rec.Body.String(); !strings.Contains(body, "at least 8 characters") {
		t.Fatalf("expected error message, got: %s", body)
	}

	// Happy path: form re-render (empty fields) with inline success message,
	// and the new password actually works.
	rec = post(testPassword, "newpass123", "newpass123")
	body := rec.Body.String()
	for _, want := range []string{
		`id="password-form"`,
		"Password updated successfully",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("password change response missing %q", want)
		}
	}

	record, err := app.FindAuthRecordByEmail("users", testEmail)
	if err != nil {
		t.Fatal(err)
	}
	if !record.ValidatePassword("newpass123") {
		t.Error("expected the new password to be active")
	}
}

