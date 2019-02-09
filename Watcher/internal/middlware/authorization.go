package middlware

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/nielsvanm/firewatch/internal/models"
)

var authRedirectURL = "/auth/login/"

// AuthorizationMiddleware is responsible for authenticating user requests
func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Info("Authorization attempt")

		// Get token from cookies
		token, err := r.Cookie("session-token")
		if err != nil {
			log.Info("Cookie missing from authorization attempt")
			http.Redirect(w, r, authRedirectURL, http.StatusSeeOther)
			return
		}

		// Get and verify the session
		sess := models.GetSessionByToken(token.Value)

		if sess.Verify() {
			sess.Save()
			next.ServeHTTP(w, r)
			return
		}

		log.Info("Failed authorization attempt")
		http.Redirect(w, r, authRedirectURL, http.StatusSeeOther)
	})
}
