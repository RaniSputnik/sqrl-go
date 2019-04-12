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

func (d *inmemoryDelegate) Authenticated(ctx context.Context, id sqrl.Identity) error {
	return nil
}
