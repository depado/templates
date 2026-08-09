package pages

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/pocketbase/pocketbase/core"

	"{{.gitserver}}/{{.owner}}/{{.name}}/models"
	"{{.gitserver}}/{{.owner}}/{{.name}}/views/layout"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestDashboardFragmentAttributes(t *testing.T) {
	props := layout.AppLayoutProps{Title: "Dashboard", CurrentPath: "/"}
	body := render(t, DashboardFragment(props))

	for _, want := range []string{
		`<title id="page-title">`,
		`id="page-heading"`,
		`id="main-content"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("dashboard fragment missing %q", want)
		}
	}
	if strings.Contains(body, `id="main-nav"`) {
		t.Error("dashboard fragment should not re-render the sidebar nav")
	}
}

func TestSettingsFragmentForms(t *testing.T) {
	props := layout.AppLayoutProps{Title: "Settings", CurrentPath: "/settings"}
	rec := core.NewRecord(&core.Collection{})
	rec.Set("name", "Test User")
	rec.Set("email", "test@example.com")
	data := SettingsData{User: &models.User{Record: rec}}
	body := render(t, SettingsFragment(props, data))

	for _, want := range []string{
		`id="password-form"`,
		`data-on:submit="@post('/settings/password', {contentType: 'form'})"`,
		`id="pw-message"`,
		`name="current_password"`,
		`name="new_password"`,
		`name="confirm_password"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("settings fragment missing %q", want)
		}
	}
}

func TestPasswordMessageKinds(t *testing.T) {
	errBody := render(t, PasswordMessage(MessageError, "nope"))
	if !strings.Contains(errBody, `role="alert"`) || !strings.Contains(errBody, "text-destructive") {
		t.Errorf("error message render wrong: %s", errBody)
	}
	okBody := render(t, PasswordMessage(MessageSuccess, "done"))
	if !strings.Contains(okBody, `role="status"`) || !strings.Contains(okBody, "text-primary") {
		t.Errorf("success message render wrong: %s", okBody)
	}
}
