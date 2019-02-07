package tools

import (
	"net/http"
	"strings"
)

// GetIPAddress takes an request and extracts the IP address without the port
func GetIPAddress(r *http.Request) string {
	fullAddr := r.RemoteAddr

	seperated := strings.Split(fullAddr, ":")

	return seperated[0]
}
