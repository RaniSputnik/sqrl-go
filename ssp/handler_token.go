package ssp

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

type User struct {
	Id string
}

// TODO: Where's the best place to put this? Extend the existing Store interface?
// Or a new store interface entirely?
type UserStore interface {
	GetAuthenticatedUser(ctx context.Context, token string) (*User, error)
}

type todoUserStore struct{}

func (s *todoUserStore) GetAuthenticatedUser(ctx context.Context, token string) (*User, error) {
	return nil, errors.New("not implemented")
}

// TokenHandler is an endpoint repsonsible for validating and exchanging the token
// issued to the client for user details so that the resource server can associate
// that SQRL user with their own copy of the user identity.
func TokenHandler(s *sqrl.Server, store UserStore, logger *log.Logger) http.Handler {
	type tokenResponse struct {
		User string `json:"user"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if err := validateToken(token); err != nil {
			logger.Printf("Token validation failed: %+v", err)
			// TODO: Standardise error responses
			// An invalid token is the same as a token that does not exist
			w.WriteHeader(http.StatusNotFound)
			return
		}

		user, err := store.GetAuthenticatedUser(r.Context(), token)
		if err != nil {
			logger.Printf("Failed to GetAuthenticatedUser: %+v", err)
			// TODO: Standardise error responses
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res := tokenResponse{User: user.Id}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			logger.Printf("User write unsuccessful: %v", err)
		}
	})
}

func validateToken(token string) error {
	// TODO: Check token was signed by this server
	// TODO: Check token has not expired

	return errors.New("not implemented")
}
