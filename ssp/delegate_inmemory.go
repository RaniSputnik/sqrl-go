package ssp

import (
	"context"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

type inmemoryDelegate struct{}

func TODODelegate() Delegate {
	return &inmemoryDelegate{}
}

func (d *inmemoryDelegate) Known(ctx context.Context, id sqrl.Identity) (bool, error) {
	return false, nil
}

func (d *inmemoryDelegate) GetSession(ctx context.Context, nut sqrl.Nut) (sqrl.Identity, SessionState, error) {
	return "", SessionUnknown, nil
}

func (d *inmemoryDelegate) Queried(ctx context.Context, id sqrl.Identity, nut sqrl.Nut) error {
	return nil
}

func (d *inmemoryDelegate) Verified(ctx context.Context, id sqrl.Identity) error {
	return nil
}

func (d *inmemoryDelegate) Redirected(ctx context.Context, id sqrl.Identity, token string) error {
	return nil
}
