package ssp_test

import (
	"context"

	sqrl "github.com/RaniSputnik/sqrl-go"
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
		GetSession struct {
			CalledWith struct {
				Ctx context.Context
				Nut sqrl.Nut
			}
			Returns struct {
				Id    sqrl.Identity
				Token string
				Err   error
			}
		}
		Queried struct {
			CalledWith struct {
				Ctx context.Context
				Id  sqrl.Identity
				Nut sqrl.Nut
			}
			Returns struct {
				Err error
			}
		}
		Verified struct {
			CalledWith struct {
				Ctx   context.Context
				Id    sqrl.Identity
				Token string
			}
			Returns struct {
				Err error
			}
		}
		Redirected struct {
			CalledWith struct {
				Ctx context.Context
				Id  sqrl.Identity
			}
			Returns struct {
				Err error
			}
		}
	}
}

func (m *mockDelegate) Known(ctx context.Context, id sqrl.Identity) (bool, error) {
	m.Func.Known.CalledWith.Ctx = ctx
	m.Func.Known.CalledWith.Id = id
	return m.Func.Known.Returns.IsKnown, m.Func.Known.Returns.Err
}

func (m *mockDelegate) GetSession(ctx context.Context, nut sqrl.Nut) (sqrl.Identity, string, error) {
	f := m.Func.GetSession
	f.CalledWith.Ctx = ctx
	f.CalledWith.Nut = nut
	return f.Returns.Id, f.Returns.Token, f.Returns.Err
}

func (m *mockDelegate) Queried(ctx context.Context, id sqrl.Identity, nut sqrl.Nut) error {
	m.Func.Queried.CalledWith.Ctx = ctx
	m.Func.Queried.CalledWith.Id = id
	m.Func.Queried.CalledWith.Nut = nut
	return m.Func.Queried.Returns.Err
}

func (m *mockDelegate) Verified(ctx context.Context, id sqrl.Identity, token string) error {
	m.Func.Verified.CalledWith.Ctx = ctx
	m.Func.Verified.CalledWith.Id = id
	m.Func.Verified.CalledWith.Token = token
	return m.Func.Verified.Returns.Err
}

func (m *mockDelegate) Redirected(ctx context.Context, id sqrl.Identity) error {
	m.Func.Redirected.CalledWith.Ctx = ctx
	m.Func.Redirected.CalledWith.Id = id
	return m.Func.Redirected.Returns.Err
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
