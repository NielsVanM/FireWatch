package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// HTTPLogMiddleware logs the HTTP request
func HTTPLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		go log.WithFields(log.Fields{
			"RemoteIP":   r.RemoteAddr,
			"URL":        r.RequestURI,
			"Method":     r.Method,
			"User-Agent": r.UserAgent(),
		}).Info("HTTP Request")

		// Continue
		next.ServeHTTP(w, r)
	})
}
