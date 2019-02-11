package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// TimingWarning is the duration where the log level changes to warning
var TimingWarning, _ = time.ParseDuration("1000ms")

// TimingError is the duration where the log level changes to error
var TimingError, _ = time.ParseDuration("5000ms")

// HTTPLogMiddleware logs the HTTP request
func HTTPLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var startTime = time.Now()

		// Handle request first
		next.ServeHTTP(w, r)

		var duration = time.Since(startTime)
		go logHTTP(r, duration)

	})
}

func logHTTP(r *http.Request, duration time.Duration) {
	// Log request
	entry := log.WithFields(log.Fields{
		"RemoteIP":   r.RemoteAddr,
		"URL":        r.RequestURI,
		"Method":     r.Method,
		"User-Agent": r.UserAgent(),
		"Duration":   duration,
	})

	// Log http rewquest on different levels
	switch time := duration; {
	case time > TimingError:
		entry.Error("HTTP Request Error Took too long")
		break
	case time > TimingWarning:
		entry.Warning("HTTP Request Warning Took too long")
		break
	default:
		entry.Info("HTTP Request")
		break
	}
}
