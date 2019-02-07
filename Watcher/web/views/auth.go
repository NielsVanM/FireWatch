package views

import (
	"fmt"
	"net/http"

	"github.com/nielsvanm/firewatch/internal/models"
	"github.com/nielsvanm/firewatch/internal/page"
	"github.com/nielsvanm/firewatch/internal/tools"
)

// LoginView allows the user to login
func LoginView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := page.NewPage("auth/login.html")
		p.Render(w)
		return
	}

	if r.Method == "POST" {
		// Get values
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Retrieve user and verify password
		u := models.GetAccountByUsername(username)
		if u == nil {
			w.Write([]byte("Invalid username/password"))
			return
		}

		success, err := u.VerifyPassword(password)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Return if the password verification failed
		if !success {
			w.Write([]byte("Invalid username/password"))
			return
		}

		// Generate session & cookie
		sess, _ := u.NewSession(tools.GetIPAddress(r))
		sess.Save()

		sessionCookie := http.Cookie{
			Name:  "session-token",
			Value: sess.SessionToken,
			Path:  "/",
		}

		// Set cookie and redirect
		http.SetCookie(w, &sessionCookie)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}
