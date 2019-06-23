package ssp

import (
	"net/http"
)

type ServerToServerAuthValidationFunc func(r *http.Request) error

type middleware func(http.Handler) http.Handler

func (server *Server) serverToServerAuthMiddleware(validator ServerToServerAuthValidationFunc) middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := validator(r); err != nil {
				server.logger.Printf("%s %s Invalid: %+v", r.Method, r.URL, err)
				// TODO: Standardise error responses
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
