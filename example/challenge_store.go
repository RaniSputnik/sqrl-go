package main

import (
	"errors"
	"sync"
)

type LoginState int8

const (
	LoginStateUnissued = LoginState(iota)
	LoginStateIssued
	LoginStateAuthenticated
	LoginStateRejected
)

func (s LoginState) String() string {
	switch s {
	case LoginStateUnissued:
		return "unissued"
	case LoginStateIssued:
		return "issued"
	case LoginStateAuthenticated:
		return "authenticated"
	case LoginStateRejected:
		return "rejected"
	}
	return "unknown"
}

var ErrInvalidState = errors.New("invalid state")

type challengeStore struct {
	vals map[string]LoginState
	sync.Mutex
}

func (s *challengeStore) Set(challenge string, nextState LoginState) error {
	if nextState == LoginStateUnissued {
		return ErrInvalidState // You may never set the state to unissued
	}

	s.Lock()
	defer s.Unlock()

	currentState, issued := s.vals[challenge]
	if !issued && nextState != LoginStateIssued {
		return ErrInvalidState // If unissued, the only valid next state is issued
	}
	if currentState == LoginStateAuthenticated ||
		currentState == LoginStateRejected {
		return ErrInvalidState // One accepted or rejected the state can not be changed
	}
	if currentState == LoginStateIssued &&
		!(nextState == LoginStateAuthenticated ||
			nextState == LoginStateRejected) {
		return ErrInvalidState // If issued, then the only valid next states are either auth or reject
	}

	s.vals[challenge] = nextState
	return nil
}

func (s *challengeStore) Get(challenge string) LoginState {
	s.Lock()
	defer s.Unlock()
	currentState, issued := s.vals[challenge]
	if !issued {
		return LoginStateUnissued
	}
	return currentState
}
