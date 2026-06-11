package router

import (
	"{{.gitserver}}/{{.owner}}/{{.name}}/models"
	"{{.gitserver}}/{{.owner}}/{{.name}}/views/layout"
	"{{.gitserver}}/{{.owner}}/{{.name}}/views/pages"
	"github.com/pocketbase/pocketbase/core"
)

func dashboardHandler(e *core.RequestEvent) error {
	h := NewHandler(e)

	props := layout.AppLayoutProps{
		Title:       "Dashboard",
		CurrentPath: "/",
		User:        models.NewUser(e.Auth),
	}

	if h.RenderPartial() {
		return render(e, pages.DashboardContent())
	}
	return render(e, pages.DashboardPage(props))
}

func settingsGetHandler(e *core.RequestEvent) error {
	h := NewHandler(e)

	props := layout.AppLayoutProps{
		Title:       "Settings",
		CurrentPath: "/settings",
		User:        models.NewUser(e.Auth),
	}

	data := pages.SettingsData{
		User: models.NewUser(e.Auth),
	}

	if h.RenderPartial() {
		return render(e, pages.SettingsContent(data))
	}
	return render(e, pages.SettingsPage(props, data))
}

func settingsPasswordHandler(e *core.RequestEvent) error {
	currentPassword := e.Request.PostFormValue("current_password")
	newPassword := e.Request.PostFormValue("new_password")
	confirmPassword := e.Request.PostFormValue("confirm_password")

	data := pages.SettingsData{
		User: models.NewUser(e.Auth),
	}

	if !e.Auth.ValidatePassword(currentPassword) {
		data.Error = "Current password is incorrect"
		return render(e, pages.SettingsContent(data))
	}
	if len(newPassword) < 8 {
		data.Error = "Password must be at least 8 characters"
		return render(e, pages.SettingsContent(data))
	}
	if newPassword != confirmPassword {
		data.Error = "Passwords do not match"
		return render(e, pages.SettingsContent(data))
	}

	e.Auth.SetPassword(newPassword)
	if err := e.App.Save(e.Auth); err != nil {
		data.Error = "Failed to update password"
		return render(e, pages.SettingsContent(data))
	}

	if err := setAuthCookie(e, e.Auth); err != nil {
		data.Error = "Failed to refresh auth"
		return render(e, pages.SettingsContent(data))
	}

	data.Success = "Password updated successfully"
	return render(e, pages.SettingsContent(data))
}
