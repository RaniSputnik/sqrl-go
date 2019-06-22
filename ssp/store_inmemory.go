package ssp

import (
	"context"
	"sync"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

type inmemoryStore struct {
	transactions      map[sqrl.Nut]*Transaction
	firstTransactions map[sqrl.Nut]sqrl.Nut
	tokens            map[sqrl.Nut]string

	sync.Mutex
}

func NewMemoryStore() Store {
	return &inmemoryStore{
		transactions:      map[sqrl.Nut]*Transaction{},
		firstTransactions: map[sqrl.Nut]sqrl.Nut{},
		tokens:            map[sqrl.Nut]string{},
	}
}

func (s *inmemoryStore) GetFirstTransaction(ctx context.Context, nut sqrl.Nut) (*Transaction, error) {
	s.Lock()
	defer s.Unlock()
	firstTransactionId, exists := s.firstTransactions[nut]
	if !exists {
		return nil, nil
	}
	return s.transactions[firstTransactionId], nil
}

func (s *inmemoryStore) SaveTransaction(ctx context.Context, t *Transaction) error {
	firstTransaction, err := s.GetFirstTransaction(ctx, t.Id)
	if err != nil {
		return err
	}
	if firstTransaction == nil {
		firstTransaction = t
	}

	s.Lock()
	defer s.Unlock()

	s.transactions[t.Id] = t
	s.firstTransactions[t.Next] = firstTransaction.Id
	return nil
}

func (s *inmemoryStore) SaveIdentSuccess(ctx context.Context, nut sqrl.Nut, token string) error {
	s.Lock()
	defer s.Unlock()
	s.tokens[nut] = token
	return nil
}

func (s *inmemoryStore) GetIdentSuccess(ctx context.Context, nut sqrl.Nut) (token string, err error) {
	s.Lock()
	defer s.Unlock()

	return s.tokens[nut], nil
}

func (s *inmemoryStore) GetIsKnown(ctx context.Context, id sqrl.Identity) (bool, error) {
	// TODO: what point is an identity is considered "known"
	// eg. is it after a successful query? Or is "known" only
	// saved after a successful ident?
	return false, nil
}
