package ssp

import (
	"context"

	sqrl "github.com/RaniSputnik/sqrl-go"
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

	GetSession(ctx context.Context, nut sqrl.Nut) (sqrl.Identity, string, error)

	// Lifecycle callbacks

	// Queried is called when a new session is started with the server.
	Queried(ctx context.Context, id sqrl.Identity, nut sqrl.Nut) error

	// Verified is called to indicate the client has been authenticated successfully.
	Verified(ctx context.Context, id sqrl.Identity, token string) error

	// Redirected is called once the redirect has been issued to the users browser
	// This may happen without 'Verified' ever being called. In that case the
	// user has been verified implicitly.
	Redirected(ctx context.Context, id sqrl.Identity) error

	// Failed is called when a login fails, should be used for auditing purposes.
	// Faild(ctx, idk, reason)
}
