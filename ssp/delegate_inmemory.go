package ssp

import (
	"context"
	"log"
	"time"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

type attempt struct {
	ID    string
	IDK   sqrl.Identity
	Nut   sqrl.Nut
	Token string

	// TODO: Failure reason?

	StartedAt time.Time
	UpdatedAt time.Time
}

type inmemoryDelegate struct {
	attempts []*attempt
}

type attemptPredicate func(*attempt) bool

func (d *inmemoryDelegate) attemptWhere(p attemptPredicate) *attempt {
	for _, a := range d.attempts {
		if p(a) {
			return a
		}
	}
	return nil
}

func identityEquals(id sqrl.Identity) attemptPredicate {
	return func(a *attempt) bool {
		return a.IDK == id
	}
}

func nutEquals(nut sqrl.Nut) attemptPredicate {
	return func(a *attempt) bool {
		return a.Nut == nut
	}
}

func TODODelegate() Delegate {
	return &inmemoryDelegate{}
}

func (d *inmemoryDelegate) Known(ctx context.Context, id sqrl.Identity) (bool, error) {
	return false, nil
}

func (d *inmemoryDelegate) GetSession(ctx context.Context, nut sqrl.Nut) (sqrl.Identity, string, error) {
	log.Printf("Delegate:GetSession nut=%s", nut)
	a := d.attemptWhere(nutEquals(nut))
	if a == nil {
		return "", "", nil
	}

	if a.Token == "" {
		return a.IDK, "", nil
	}

	return a.IDK, a.Token, nil
}

func (d *inmemoryDelegate) Queried(ctx context.Context, id sqrl.Identity, nut sqrl.Nut) error {
	log.Printf("Delegate:Queried id=%s, nut=%s", id, nut)
	now := time.Now()
	newAttempt := attempt{
		IDK:       id,
		Nut:       nut,
		StartedAt: now,
		UpdatedAt: now,
	}
	d.attempts = append(d.attempts, &newAttempt)
	return nil
}

func (d *inmemoryDelegate) Verified(ctx context.Context, id sqrl.Identity, token string) error {
	log.Printf("Delegate:Verified id=%s, token=%s", id, token)
	// TODO: What do we actually want to do with this callback?
	// Would we prefer
	a := d.attemptWhere(identityEquals(id))
	if a != nil {
		a.Token = token
		a.UpdatedAt = time.Now()
	}
	return nil
}

// TODO: We might want to pass 'nut' here instead of identity
func (d *inmemoryDelegate) Redirected(ctx context.Context, id sqrl.Identity) error {
	log.Printf("Delegate:Redirected id=%s", id)
	// TODO: What do we actually want to do with this callback?
	return nil
}
