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

	SaveTransaction(ctx context.Context, t *Transaction) error
	SaveIdentSuccess(ctx context.Context, nut sqrl.Nut, token string) error

	// GetIdentSuccess returns a previously saved token for a given transaction nut
	// if such a token exists. An empty string will be returned if the given nut
	// has not yet been saved as successful.
	GetIdentSuccess(ct context.Context, nut sqrl.Nut) (token string, err error)

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

// // TODO: Save this to the DB to 'cache' a reference
// type sessionTransaction struct {
// 	// The id of the transaction that started this session
// 	// Note: If Id == Transaction this is the first transaction
// 	// in a session.
// 	Id sqrl.Nut
// 	// The id of the transaction this link refers to.
// 	Transaction sqrl.Nut
// }

// // Populated using session transactions
// // probably not needed but here for reference
// type session struct {
// 	Transactions []Transaction
// }
