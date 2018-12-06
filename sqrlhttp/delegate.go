package sqrlhttp

import (
	"context"

	"github.com/RaniSputnik/sqrl-go"
)

// Delegate is an interface that allows a service to implement
// it's own database and SQRL storage mechanisms.
type Delegate interface {
	// Known is called to determine whether the given
	// SQRL idenity has previously logged in to the server.
	// This function should check the database or cache
	// for a user with the given id.
	//
	// Return true/false to indicate if the user exists and
	// an error if the determination was unsuccessful.
	Known(ctx context.Context, id sqrl.Identity) (bool, error)

	// Authenticated is called when a client has successfully
	// identified itself with the SQRL server. This identity
	// should now be considered to be logged in.
	Authenticated(ctx context.Context, id sqrl.Identity) error
}
