package ssp

import (
	"encoding/json"
	"log"
	"net/http"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

// TokenHandler is an endpoint repsonsible for validating and exchanging the token
// issued to the client for user details so that the resource server can associate
// that SQRL user with their own copy of the user identity.
func TokenHandler(s *sqrl.Server, tokens *TokenGenerator, logger *log.Logger) http.Handler {
	type tokenResponse struct {
		User string `json:"user"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		userId, err := tokens.ValidateToken(token)
		if err != nil {
			logger.Printf("Token validation failed: %+v", err)
			// TODO: Standardise error responses
			// An invalid token is the same as a token that does not exist
			w.WriteHeader(http.StatusNotFound)
			return
		}

		res := tokenResponse{User: userId}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			logger.Printf("User write unsuccessful: %v", err)
		}
	})
}
