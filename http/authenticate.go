package http

import (
	"log"
	"net/http"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

const xFormURLEncoded = "application/x-www-form-urlencoded"

func Authenticate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.UserAgent() != sqrl.V1 || r.Header.Get("Content-Type") != xFormURLEncoded {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("Failed to parse form: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}
