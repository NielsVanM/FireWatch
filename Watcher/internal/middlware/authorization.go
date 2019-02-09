package middlware

import (
	"fmt"
	"net/http"

	"github.com/nielsvanm/firewatch/internal/models"
)

var authRedirectURL = "/auth/login/"

// AuthorizationMiddleware is responsible for authenticating user requests
func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get token from cookies
		token, err := r.Cookie("session-token")
		if err != nil {
			http.Redirect(w, r, authRedirectURL, http.StatusSeeOther)
			fmt.Println("Authorization", err.Error())
			return
		}

		// Get and verify the session
		sess := models.GetSessionByToken(token.Value)

		if sess.Verify() {
			sess.Save()
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, authRedirectURL, http.StatusSeeOther)
	})
}
