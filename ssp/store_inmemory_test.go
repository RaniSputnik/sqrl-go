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
		_ = s.SaveTransaction(ctx, &ssp.Transaction{
			Id:   nut,
			Next: sqrl.Nut("someothernut"),
		})
		transaction, _ := s.GetFirstTransaction(ctx, nut)
		assert.Nil(t, transaction)
	})

	t.Run("ReturnFirstTransactionIfExists", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		firstNut := sqrl.Nut("neverusedbefore")
		thisNut := sqrl.Nut("someothernut")
		_ = s.SaveTransaction(ctx, &ssp.Transaction{
			Id:   firstNut,
			Next: thisNut,
		})

		transaction, _ := s.GetFirstTransaction(ctx, thisNut)
		if assert.NotNil(t, transaction) {
			assert.Equal(t, firstNut, transaction.Id)
		}
	})

	t.Run("ReturnsFirstTransactionInALongSession", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		firstNut := sqrl.Nut("firstnut")
		secondNut := sqrl.Nut("secondnut")
		thirdNut := sqrl.Nut("thirdnut")
		_ = s.SaveTransaction(ctx, &ssp.Transaction{
			Id:   firstNut,
			Next: secondNut,
		})
		_ = s.SaveTransaction(ctx, &ssp.Transaction{
			Id:   secondNut,
			Next: thirdNut,
		})
		_ = s.SaveTransaction(ctx, &ssp.Transaction{
			Id:   thirdNut,
			Next: sqrl.Nut("someothernut"),
		})

		transaction, _ := s.GetFirstTransaction(ctx, thirdNut)
		if assert.NotNil(t, transaction) {
			assert.Equal(t, firstNut, transaction.Id)
		}
	})
}

func TestMemoryStoreIdent(t *testing.T) {
	ctx := context.TODO()
	knownNut := sqrl.Nut("somenut")

	t.Run("ReturnsPreviouslySavedToken", func(t *testing.T) {
		s := ssp.NewMemoryStore()
		givenToken := "abcdef1234567890"

		_ = s.SaveIdentSuccess(ctx, knownNut, givenToken)
		gotToken, err := s.GetIdentSuccess(ctx, knownNut)

		assert.Nil(t, err)
		assert.Equal(t, givenToken, gotToken)
	})
}

func TestMemoryStoreUsers(t *testing.T) {
	// TODO: Add tests for user functions
}
