package views

import (
	"fmt"
	"net/http"

	"github.com/nielsvanm/firewatch/internal/models"
	"github.com/nielsvanm/firewatch/internal/page"

	log "github.com/sirupsen/logrus"
)

// LoginView allows the user to login
func LoginView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := page.NewPage("auth/login.html")

		// Get referal
		refURL, ok := r.URL.Query()["nextPage"]
		log.Info(refURL)
		if ok {
			p.AddContext("nextPage", refURL[0])
			p.AddContext("message", page.NewMessage(page.MessageInfo, "You need to be logged in to access this page.", false))
		}

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
			invalidUsernamePassword(w)
			return
		}

		success, _ := u.VerifyPassword(password)

		// Return if the password verification failed
		if !success {
			invalidUsernamePassword(w)
			return
		}

		// Generate session & cookie
		sess, _ := u.NewSession()
		sess.Save()

		sessionCookie := http.Cookie{
			Name:  "session-token",
			Value: sess.SessionToken,
			Path:  "/",
		}

		// Set cookie and redirect
		http.SetCookie(w, &sessionCookie)

		// Check if there is a referal
		if nextPage := r.URL.Query()["nextPage"]; nextPage != nil {
			http.Redirect(w, r, nextPage[0], http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// invalidUsernamePassword writes the login page with an invalid username/password
// message
func invalidUsernamePassword(w http.ResponseWriter) {
	p := page.NewPage("auth/login.html")
	p.AddContext("message", page.NewMessage(
		page.MessageWarning,
		"Invalid Username/Password",
		false,
	))
	p.Render(w)
}

// LogoutView deletes the session of the user when he visits the view, after
// deletion of the session it redirects the user to the login page.
func LogoutView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ck, err := r.Cookie("session-token")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		sess := models.GetSessionByToken(ck.Value)
		sess.Delete()

		http.Redirect(w, r, "/auth/login/", http.StatusSeeOther)
	}

	if r.Method == "GET" {
		http.Redirect(w, r, "/auth/login/", http.StatusSeeOther)
	}
}
