package ssp

import (
	"context"
	"errors"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

type Store interface {
	// GetFirstTransaction returns the transaction that started an exchange between
	// a SQRL client and SSP server. If no error or transaction is returned then
	// the current transaction is the first transaction in the exchange.
	GetFirstTransaction(ctx context.Context, nut sqrl.Nut) (*Transaction, error)

	// SaveTransaction stores a verified transaction in the DB.
	SaveTransaction(ctx context.Context, t *Transaction) error

	// SaveIdentSuccess stores a successful ident query from a client. The token
	// that will be returned to the client is stored to allow for retrieval
	// (for the pag.sqrl endpoint).
	SaveIdentSuccess(ctx context.Context, nut sqrl.Nut, token string) error

	// GetIdentSuccess returns a previously saved token for a given transaction nut
	// if such a token exists. An empty string will be returned if the given nut
	// has not yet been saved as successful.
	GetIdentSuccess(ctx context.Context, nut sqrl.Nut) (token string, err error)

	// Superceeded by UserStore.GetUserForToken
	// GetIsKnown(ctx context.Context, id sqrl.Identity) (bool, error)
}

// TODO: Save this to the DB each time a new query is made
// We can use this to tie nuts back to their OG id.
type Transaction struct {
	Id   sqrl.Nut
	Next sqrl.Nut
	// Client string
	// Server string
}

// TODO: We probably don't need this
// We should instead store the userId in token itself
// Along with the token expiry
type TokenStore interface {
	// This method is pretty gross - it'd be easy to mix up token and user id
	// TODO: Stronger type for either token, userId or both
	SaveToken(ctx context.Context, token string, userId string) error
	GetUserForToken(ctx context.Context, token string) (userId string, err error)
}

type todoTokenStore struct{}

func (s *todoTokenStore) SaveToken(ctx context.Context, token string, userId string) error {
	return errors.New("not implemented")
}

func (s *todoTokenStore) GetUserForToken(ctx context.Context, token string) (userId string, err error) {
	return "", errors.New("not implemented")
}

type UserStore interface {
	CreateUser(ctx context.Context, idk sqrl.Identity) (*User, error)

	// GetByIdentity returns a user from the given identity key.
	// If no user is found, a nil user will be returned with no error.
	// TODO: Clarify exactly when a user should be saved
	// is it after a successful query? Or after successful ident?
	// see: https://github.com/RaniSputnik/sqrl-go/issues/25
	GetUserByIdentity(ctx context.Context, idk sqrl.Identity) (*User, error)

	// TODO: Get by previous identities
}

type User struct {
	Id  string
	Idk sqrl.Identity
	// TODO: Do we need to store previous identity keys?
}

type todoUserStore struct{}

func (s *todoUserStore) CreateUser(ctx context.Context, idk sqrl.Identity) (*User, error) {
	return nil, errors.New("not implemented")
}

func (s *todoUserStore) GetUserByIdentity(ctx context.Context, idk sqrl.Identity) (*User, error) {
	return nil, errors.New("not implemented")
}
