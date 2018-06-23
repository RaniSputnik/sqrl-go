package http

import (
	"net/http"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

func Authenticate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.UserAgent() != sqrl.V1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}
