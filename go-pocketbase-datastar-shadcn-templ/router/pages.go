package router

import (
	"github.com/pocketbase/pocketbase/core"

	"{{.gitserver}}/{{.owner}}/{{.name}}/models"
	"{{.gitserver}}/{{.owner}}/{{.name}}/views/layout"
	"{{.gitserver}}/{{.owner}}/{{.name}}/views/pages"
)

func appLayoutProps(title, currentPath string, e *core.RequestEvent) layout.AppLayoutProps {
	return layout.AppLayoutProps{
		Title:            title,
		CurrentPath:      currentPath,
		User:             models.NewUser(e.Auth),
		SidebarCollapsed: sidebarCollapsed(e),
	}
}

// sidebarCollapsed restores the sidebar state persisted by the component
// script ("sidebar_state" cookie; "false" means collapsed, like shadcn's
// useSidebar side-effect).
func sidebarCollapsed(e *core.RequestEvent) bool {
	c, err := e.Request.Cookie("sidebar_state")
	return err == nil && c.Value == "false"
}

func dashboardHandler(e *core.RequestEvent) error {
	props := appLayoutProps("Dashboard", "/", e)

	if isDatastar(e) {
		return render(e, pages.DashboardFragment(props))
	}
	return render(e, pages.DashboardPage(props))
}

func settingsGetHandler(e *core.RequestEvent) error {
	props := appLayoutProps("Settings", "/settings", e)
	data := pages.SettingsData{User: models.NewUser(e.Auth)}

	if isDatastar(e) {
		return render(e, pages.SettingsFragment(props, data))
	}
	return render(e, pages.SettingsPage(props, data))
}

func settingsPasswordHandler(e *core.RequestEvent) error {
	currentPassword := e.Request.PostFormValue("current_password")
	newPassword := e.Request.PostFormValue("new_password")
	confirmPassword := e.Request.PostFormValue("confirm_password")

	sse := newSSE(e)

	if !e.Auth.ValidatePassword(currentPassword) {
		return sse.PatchElementTempl(pages.PasswordMessage(pages.MessageError, "Current password is incorrect"))
	}
	if len(newPassword) < 8 {
		return sse.PatchElementTempl(pages.PasswordMessage(pages.MessageError, "Password must be at least 8 characters"))
	}
	if newPassword != confirmPassword {
		return sse.PatchElementTempl(pages.PasswordMessage(pages.MessageError, "Passwords do not match"))
	}

	e.Auth.SetPassword(newPassword)
	if err := e.App.Save(e.Auth); err != nil {
		return sse.PatchElementTempl(pages.PasswordMessage(pages.MessageError, "Failed to update password"))
	}

	// The token is derived from the password hash; reissue it after a change.
	if err := setAuthCookie(e, e.Auth); err != nil {
		return sse.PatchElementTempl(pages.PasswordMessage(pages.MessageError, "Failed to refresh auth"))
	}

	// Success: re-render the form (clears the password fields) with the
	// inline success message.
	return sse.PatchElementTempl(pages.PasswordForm("Password updated successfully"))
}
