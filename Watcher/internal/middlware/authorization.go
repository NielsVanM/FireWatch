package middlware

import (
	"fmt"
	"net/http"

	"github.com/nielsvanm/firewatch/internal/models"
	"github.com/nielsvanm/firewatch/internal/tools"
)

// AuthorizationMiddleware is responsible for authenticating user requests
func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get token from cookies
		token, err := r.Cookie("session-token")
		if err != nil {
			http.Redirect(w, r, "/auth/login/", http.StatusSeeOther)
			fmt.Println("Authorization", err.Error())
			return
		}

		// Get and verify the session
		sess := models.GetSessionByToken(token.Value)
		if sess.Verify(tools.GetIPAddress(r)) {
			sess.Save()
			next.ServeHTTP(w, r)
			return
		}
		w.Write([]byte("User not authenticated or session expired. Please log in"))
	})
}
