package ssp_test

import (
	"testing"
	"time"

	"github.com/RaniSputnik/sqrl-go/ssp"
	"github.com/stretchr/testify/assert"

	"encoding/base64"
)

var anyKey = make([]byte, 16)

func TestTokenGeneration(t *testing.T) {
	generator := ssp.DefaultExchange(anyKey, time.Minute)

	t.Run("ReturnsAToken", func(t *testing.T) {
		got := generator.Token("someUser")
		t.Logf("Generated token: %s", got)
		assert.NotEmpty(t, got)
	})

	t.Run("Base64EncodesTheToken", func(t *testing.T) {
		got := generator.Token("someUser")
		_, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(string(got))
		t.Logf("Generated token: %s", got)
		assert.Nil(t, err)
	})

	t.Run("EncodesAReallyLongUserId", func(t *testing.T) {
		got := generator.Token("areallylonguseridthatiscompletelyunreasonable")
		t.Logf("Generated token: %s", got)
		assert.NotEmpty(t, got)
	})
}

func TestTokenValidation(t *testing.T) {
	generator := ssp.DefaultExchange(anyKey, time.Minute)
	validToken := generator.Token("someUser")

	t.Run("ReturnsNoErrorForValidToken", func(t *testing.T) {
		_, err := generator.Validate(validToken)
		assert.Nil(t, err)
	})

	t.Run("ReturnsUserIDForValidToken", func(t *testing.T) {
		expectedUserID := "someUser"
		uid, _ := generator.Validate(validToken)
		assert.Equal(t, expectedUserID, uid)
	})

	t.Run("ReturnsTokenExpiredForOldToken", func(t *testing.T) {
		veryShortExpiry := time.Millisecond
		exchangeWithShortExpiry := ssp.DefaultExchange(anyKey, veryShortExpiry)
		expiredToken := exchangeWithShortExpiry.Token("someUser")

		time.Sleep(veryShortExpiry * 3)

		_, err := exchangeWithShortExpiry.Validate(expiredToken)
		assert.Equal(t, ssp.ErrTokenExpired, err)
	})

	t.Run("ReturnsTokenFormatInvalidForEmptyToken", func(t *testing.T) {
		_, err := generator.Validate("")
		assert.Equal(t, ssp.ErrTokenFormatInvalid, err)
	})

	t.Run("ReturnsTokenFormatInvalidForRandomString", func(t *testing.T) {
		_, err := generator.Validate("someinvalidtoken")
		assert.Equal(t, ssp.ErrTokenFormatInvalid, err)
	})
}
