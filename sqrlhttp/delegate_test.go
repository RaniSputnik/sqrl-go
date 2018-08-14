package sqrlhttp_test

import (
	"context"

	"github.com/RaniSputnik/sqrl-go"
)

type mockDelegate struct {
	Func struct {
		Known struct {
			CalledWith struct {
				Ctx context.Context
				Id  sqrl.Identity
			}
			Returns struct {
				IsKnown bool
				Err     error
			}
		}
	}
}

func (m *mockDelegate) Known(ctx context.Context, id sqrl.Identity) (bool, error) {
	m.Func.Known.CalledWith.Ctx = ctx
	m.Func.Known.CalledWith.Id = id
	return m.Func.Known.Returns.IsKnown, m.Func.Known.Returns.Err
}

// Language helpers

func NewDelegate() *mockDelegate {
	return &mockDelegate{}
}

func (m *mockDelegate) ReturnsUnknownIdentity() *mockDelegate {
	m.Func.Known.Returns.IsKnown = false
	m.Func.Known.Returns.Err = nil
	return m
}

func (m *mockDelegate) ReturnsKnownIdentity() *mockDelegate {
	m.Func.Known.Returns.IsKnown = true
	m.Func.Known.Returns.Err = nil
	return m
}
