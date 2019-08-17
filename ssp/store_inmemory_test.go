package ssp_test

import (
	"context"
	"testing"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/RaniSputnik/sqrl-go/ssp"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStoreIsCreated(t *testing.T) {
	s := ssp.NewMemoryStore()
	assert.NotNil(t, s)
}

func TestMemoryStoreGetFirstTransaction(t *testing.T) {
	ctx := context.TODO()

	t.Run("ReturnsNoTransactionWhenNutNotFound", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		transaction, _ := s.GetFirstTransaction(ctx, sqrl.Nut("neverusedbefore"))
		assert.Nil(t, transaction)
	})

	t.Run("ReturnsNoTransactionIfExistsButIsFirstTransaction", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		nut := sqrl.Nut("neverusedbefore")
		_ = s.SaveTransaction(ctx, &sqrl.Transaction{
			Request: &sqrl.Request{
				Nut:      nut,
				Client:   "some-client",
				Server:   "some-server",
				Ids:      "some-signature",
				ClientIP: "10.0.0.1",
			},
			Next: sqrl.Nut("someothernut"),
		})
		transaction, _ := s.GetFirstTransaction(ctx, nut)
		assert.Nil(t, transaction)
	})

	t.Run("ReturnFirstTransactionIfExists", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		firstNut := sqrl.Nut("neverusedbefore")
		thisNut := sqrl.Nut("someothernut")
		_ = s.SaveTransaction(ctx, &sqrl.Transaction{
			Request: &sqrl.Request{
				Nut:      firstNut,
				Client:   "some-client",
				Server:   "some-server",
				Ids:      "some-signature",
				ClientIP: "10.0.0.1",
			},
			Next: thisNut,
		})

		transaction, _ := s.GetFirstTransaction(ctx, thisNut)
		if assert.NotNil(t, transaction) {
			assert.Equal(t, firstNut, transaction.Nut)
		}
	})

	t.Run("ReturnsFirstTransactionInALongSession", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		firstNut := sqrl.Nut("firstnut")
		secondNut := sqrl.Nut("secondnut")
		thirdNut := sqrl.Nut("thirdnut")
		_ = s.SaveTransaction(ctx, &sqrl.Transaction{
			Request: &sqrl.Request{Nut: firstNut},
			Next:    secondNut,
		})
		_ = s.SaveTransaction(ctx, &sqrl.Transaction{
			Request: &sqrl.Request{Nut: secondNut},
			Next:    thirdNut,
		})
		_ = s.SaveTransaction(ctx, &sqrl.Transaction{
			Request: &sqrl.Request{Nut: thirdNut},
			Next:    sqrl.Nut("someothernut"),
		})

		transaction, _ := s.GetFirstTransaction(ctx, thirdNut)
		if assert.NotNil(t, transaction) {
			assert.Equal(t, firstNut, transaction.Nut)
		}
	})
}

func TestMemoryStoreIdent(t *testing.T) {
	ctx := context.TODO()
	knownNut := sqrl.Nut("somenut")

	t.Run("ReturnsPreviouslySavedToken", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		givenToken := ssp.Token("abcdef1234567890")

		_ = s.SaveIdentSuccess(ctx, knownNut, givenToken)
		gotToken, err := s.GetIdentSuccess(ctx, knownNut)

		assert.Nil(t, err)
		assert.Equal(t, givenToken, gotToken)
	})
}

func TestMemoryStoreUsers(t *testing.T) {
	ctx := context.TODO()

	t.Run("CreatesUserReturnsAUser", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		user, err := s.CreateUser(ctx, "someidk")
		assert.Nil(t, err)
		assert.NotNil(t, user)
	})

	// TODO: What to do if we attempt to create a user with the same IDK?

	t.Run("CreateUserUsesAUniqueIDForTheUser", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		user1, _ := s.CreateUser(ctx, "idk1")
		user2, _ := s.CreateUser(ctx, "idk2")
		if assert.NotNil(t, user1) && assert.NotNil(t, user2) {
			assert.NotEmpty(t, user1.Id)
			assert.NotEmpty(t, user2.Id)
			assert.NotEqual(t, user1.Id, user2.Id)
		}
	})

	t.Run("CreateUserSetsTheNewUsersIDKCorrectly", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		idk := sqrl.Identity("someidk")
		user, _ := s.CreateUser(ctx, idk)
		if assert.NotNil(t, user) {
			assert.Equal(t, idk, user.Idk)
		}
	})

	t.Run("GetUserByIdentityReturnsKnownUser", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		idk := sqrl.Identity("someidk")
		user, _ := s.CreateUser(ctx, idk)

		fetchedUser, err := s.GetUserByIdentity(ctx, idk)
		assert.Nil(t, err)
		if assert.NotNil(t, fetchedUser) {
			assert.Equal(t, user.Id, fetchedUser.Id)
		}
	})

	t.Run("GetUserByIdentityReturnsNilIfUserNotFound", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		fetchedUser, err := s.GetUserByIdentity(ctx, "someidk")
		assert.Nil(t, err)
		assert.Nil(t, fetchedUser)
	})
}
