package ssp

import (
	"context"

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
	GetIdentSuccess(ct context.Context, nut sqrl.Nut) (token string, err error)

	// GetIsKnown returns whether the given sqrl identity has been seen in
	// previous SQRL transactions.
	// TODO: Clarify exactly when a identity is considered "known"
	// is it after a successful query? Or after successful ident?
	GetIsKnown(ctx context.Context, id sqrl.Identity) (bool, error)
}

// TODO: Save this to the DB each time a new query is made
// We can use this to tie nuts back to their OG id.
type Transaction struct {
	Id   sqrl.Nut
	Next sqrl.Nut
	// Client string
	// Server string
}
