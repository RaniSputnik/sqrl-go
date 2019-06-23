package ssp

import (
	"encoding/json"
	"net/http"
)

// TokenHandler is an endpoint repsonsible for validating and exchanging the token
// issued to the client for user details so that the resource server can associate
// that SQRL user with their own copy of the user identity.
func (server *Server) TokenHandler(tokens *TokenGenerator) http.Handler {
	type tokenResponse struct {
		User string `json:"user"`
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		userId, err := tokens.ValidateToken(token)
		if err != nil {
			server.logger.Printf("Token validation failed: %+v", err)
			// TODO: Standardise error responses
			// An invalid token is the same as a token that does not exist
			w.WriteHeader(http.StatusNotFound)
			return
		}

		res := tokenResponse{User: userId}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			server.logger.Printf("User write unsuccessful: %v", err)
		}
	})

	return server.protect(h)
}
