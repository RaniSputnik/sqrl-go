package ssp_test

import (
	"context"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/RaniSputnik/sqrl-go/ssp"
)

type mockStore struct {
	Func struct {
		GetFirstTransaction struct {
			CalledWith struct {
				Ctx context.Context
				Nut sqrl.Nut
			}
			Returns struct {
				Transaction *ssp.Transaction
				Err         error
			}
		}
		SaveTransaction struct {
			CalledWith struct {
				Ctx         context.Context
				Transaction *ssp.Transaction
			}
			Returns struct {
				Err error
			}
		}
		SaveIdentSuccess struct {
			CalledWith struct {
				Ctx   context.Context
				Nut   sqrl.Nut
				Token string
			}
			Returns struct {
				Err error
			}
		}
		GetIdentSuccess struct {
			CalledWith struct {
				Ctx context.Context
				Nut sqrl.Nut
			}
			Returns struct {
				Token string
				Err   error
			}
		}
		CreateUser struct {
			CalledWith struct {
				Ctx context.Context
				Idk sqrl.Identity
			}
			Returns struct {
				User *ssp.User
				Err  error
			}
		}
		GetUserByIdentity struct {
			CalledWith struct {
				Ctx context.Context
				Idk sqrl.Identity
			}
			Returns struct {
				User *ssp.User
				Err  error
			}
		}
	}
}

func (m *mockStore) GetFirstTransaction(ctx context.Context, nut sqrl.Nut) (*ssp.Transaction, error) {
	m.Func.GetFirstTransaction.CalledWith.Ctx = ctx
	m.Func.GetFirstTransaction.CalledWith.Nut = nut
	return m.Func.GetFirstTransaction.Returns.Transaction, m.Func.GetFirstTransaction.Returns.Err
}

func (m *mockStore) SaveTransaction(ctx context.Context, t *ssp.Transaction) error {
	m.Func.SaveTransaction.CalledWith.Ctx = ctx
	m.Func.SaveTransaction.CalledWith.Transaction = t
	return m.Func.SaveTransaction.Returns.Err
}

func (m *mockStore) SaveIdentSuccess(ctx context.Context, nut sqrl.Nut, token string) error {
	m.Func.SaveIdentSuccess.CalledWith.Ctx = ctx
	m.Func.SaveIdentSuccess.CalledWith.Nut = nut
	m.Func.SaveIdentSuccess.CalledWith.Token = token
	return m.Func.SaveIdentSuccess.Returns.Err
}

func (m *mockStore) GetIdentSuccess(ctx context.Context, nut sqrl.Nut) (token string, err error) {
	m.Func.GetIdentSuccess.CalledWith.Ctx = ctx
	m.Func.GetIdentSuccess.CalledWith.Nut = nut
	return m.Func.GetIdentSuccess.Returns.Token, m.Func.GetIdentSuccess.Returns.Err
}

func (m *mockStore) CreateUser(ctx context.Context, idk sqrl.Identity) (*ssp.User, error) {
	m.Func.CreateUser.CalledWith.Ctx = ctx
	m.Func.CreateUser.CalledWith.Idk = idk
	return m.Func.CreateUser.Returns.User, m.Func.CreateUser.Returns.Err
}

func (m *mockStore) GetUserByIdentity(ctx context.Context, idk sqrl.Identity) (*ssp.User, error) {
	m.Func.GetUserByIdentity.CalledWith.Ctx = ctx
	m.Func.GetUserByIdentity.CalledWith.Idk = idk
	return m.Func.GetUserByIdentity.Returns.User, m.Func.GetUserByIdentity.Returns.Err
}

// Language helpers

func NewStore() *mockStore {
	return &mockStore{}
}

func (m *mockStore) ReturnsUnknownIdentity() *mockStore {
	m.Func.GetUserByIdentity.Returns.User = nil
	m.Func.GetUserByIdentity.Returns.Err = nil
	return m
}

func (m *mockStore) ReturnsKnownIdentity() *mockStore {
	m.Func.GetUserByIdentity.Returns.User = &ssp.User{
		Id:  "someuser",
		Idk: "abc123",
	}
	m.Func.GetUserByIdentity.Returns.Err = nil
	return m
}
