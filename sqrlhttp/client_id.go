package sqrlhttp

import "net/http"

func clientID(r *http.Request) string {
	// TODO: X-Forwarded-For
	// TODO: Include user agent if available
	return r.RemoteAddr
}
