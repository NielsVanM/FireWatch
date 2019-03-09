package middleware

import (
	"net/http"

	"github.com/gorilla/context"

	"github.com/nielsvanm/firewatch/core/models"
)

var authRedirectURL = "/auth/login/"

// AuthorizationMiddleware is responsible for authenticating user requests
func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get token from cookies
		token, err := r.Cookie("session-token")
		if err != nil {
			http.Redirect(w, r, authRedirectURL+"?nextPage="+r.RequestURI, http.StatusSeeOther)
			return
		}

		// Get and verify the session
		sess := models.GetSessionByToken(token.Value)

		if sess.Verify() {
			sess.Save()

			// Retrieve user
			user := models.GetAccountByID(sess.UserID)
			context.Set(r, "user", user)

			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, authRedirectURL+"?nextPage="+r.RequestURI, http.StatusSeeOther)
	})
}
