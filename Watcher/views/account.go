package views

import (
	"net/http"

	"github.com/nielsvanm/firewatch/internal/tools"

	"github.com/gorilla/context"
	"github.com/nielsvanm/firewatch/internal/models"
	"github.com/nielsvanm/firewatch/internal/page"
)

// AccountView shows the account overview page
func AccountView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Handle wegpage
		p := page.NewPage("components/base.html", "account/overview.html")
		p.AddContext("user", context.Get(r, "user"))
		p.Render(w)
	}
}

// LogOutAllDevicesView logs the currently logged in user out from all devices
// he has
func LogOutAllDevicesView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user := context.Get(r, "user").(*models.Account)
		user.DeleteAllSessions()

		http.Redirect(w, r, "/account/", http.StatusSeeOther)
	}

	w.Write([]byte("Invalid Request, Need POST not " + r.Method))
}

// ChangePasswordView allows the user to changehis password
func ChangePasswordView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := page.NewPage("components/base.html", "account/change_pass.html")
		p.AddContext("user", context.Get(r, "user"))
		p.Render(w)
	}

	if r.Method == "POST" {
		// Get the values from the form
		oldPassword := r.FormValue("password_old")
		newPassword := r.FormValue("password_new")
		newPasswordDuplicate := r.FormValue("password_new_duplicate")

		user := context.Get(r, "user").(*models.Account)

		// Check old password
		passwordMatch, _ := user.VerifyPassword(oldPassword)
		if !passwordMatch {
			p := page.NewPage("components/base.html", "account/change_pass.html")
			p.AddContext("user", context.Get(r, "user"))
			p.AddContext("message", page.NewMessage(
				page.MessageWarning, "Old password is not valid.", false))
			p.Render(w)
			return
		}

		// Verify password
		if !tools.PasswordVerification(newPassword) {
			p := page.NewPage("components/base.html", "account/change_pass.html")
			p.AddContext("user", context.Get(r, "user"))
			p.AddContext("message", page.NewMessage(
				page.MessageWarning, "The new password isn't valid, it needs to be at least 8 characters long and contain a minimum of one numerical value (0-9)", false))
			p.Render(w)
			return
		}

		// Compare new passwords
		if newPassword == newPasswordDuplicate {
			user.SetPassword(newPassword)
			user.Save()

			http.Redirect(w, r, "/account/", http.StatusSeeOther)
			return
		}

		p := page.NewPage("components/base.html", "account/change_pass.html")
		p.AddContext("user", context.Get(r, "user"))
		p.AddContext("message", page.NewMessage(
			page.MessageWarning, "New passwords do not match", false))
		p.Render(w)
	}
}

// DeleteAccountView allows the end user te delete his/her account after providng
// their password again
func DeleteAccountView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Return form
		p := page.NewPage("components/base.html", "account/delete_account.html")
		p.AddContext("user", context.Get(r, "user"))
		p.Render(w)
	}

	if r.Method == "POST" {
		password := r.FormValue("password")
		user := context.Get(r, "user").(*models.Account)

		passMatch, _ := user.VerifyPassword(password)
		if passMatch {
			user.Delete()
			http.Redirect(w, r, "/auth/login/", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/account/", http.StatusSeeOther)
	}
}
