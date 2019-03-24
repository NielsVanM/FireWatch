package middleware

import (
	"net/http"

	"github.com/gorilla/context"

	"github.com/nielsvanm/firewatch/api"
	"github.com/nielsvanm/firewatch/core/models"
)

var authRedirectURL = "/auth/login/"

// AuthorizationMiddleware is responsible for authenticating user requests
func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get token from headers
		token := r.Header.Get("Authorization")
		if token == "" {
			api.NewResp(false, api.StatusInvalidToken).Write(w)
			return
		}

		// Get and verify the session
		sess := models.GetSessionByToken(token)

		if sess.Verify() {
			go sess.Save()

			// Set necessarry context
			user := models.GetAccountByID(sess.UserID)
			context.Set(r, "user", user)
			context.Set(r, "session", sess)

			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, authRedirectURL+"?nextPage="+r.RequestURI, http.StatusSeeOther)
	})
}
