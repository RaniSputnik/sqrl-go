package main

import (
	"context"

	"github.com/RaniSputnik/sqrl-go"
)

type delegate struct{}

func (d *delegate) Known(ctx context.Context, id sqrl.Identity) (bool, error) {
	// TODO: Delegate implementation
	return false, nil
}
