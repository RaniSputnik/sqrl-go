package ssp

import (
	"context"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

type TransactionStore interface {
	// GetFirstTransaction returns the transaction that started an exchange between
	// a SQRL client and SSP server. If no error or transaction is returned then
	// the current transaction is the first transaction in the exchange.
	GetFirstTransaction(ctx context.Context, nut sqrl.Nut) (*sqrl.Transaction, error)

	// SaveTransaction stores a verified transaction in the DB.
	SaveTransaction(ctx context.Context, t *sqrl.Transaction) error

	// SaveIdentSuccess stores a successful ident query from a client. The token
	// that will be returned to the client is stored to allow for retrieval
	// (for the pag.sqrl endpoint).
	SaveIdentSuccess(ctx context.Context, nut sqrl.Nut, token Token) error

	// GetIdentSuccess returns a previously saved token for a given transaction nut
	// if such a token exists. An empty string will be returned if the given nut
	// has not yet been saved as successful.
	GetIdentSuccess(ctx context.Context, nut sqrl.Nut) (token Token, err error)
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

type Store interface {
	TransactionStore
	UserStore
}
