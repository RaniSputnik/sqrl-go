package ssp

import (
	"net/http"
)

type ServerToServerAuthValidationFunc func(r *http.Request) error

func (server *Server) protect(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := server.validator(r); err != nil {
			server.logger.Printf("%s %s Invalid: %+v", r.Method, r.URL, err)
			// TODO: Standardise error responses
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func noProtection(r *http.Request) error {
	return nil
}
