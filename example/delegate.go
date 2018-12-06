package main

import (
	"context"
	"errors"

	"github.com/RaniSputnik/sqrl-go"
)

type delegate struct{}

func (d *delegate) Known(ctx context.Context, id sqrl.Identity) (bool, error) {
	// TODO: Delegate implementation
	return false, nil
}

func (d *delegate) Authenticated(ctx context.Context, id sqrl.Identity) error {
	// TODO: How do we save this state?
	return errors.New("not implemented")
}
