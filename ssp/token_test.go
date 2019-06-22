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
	generator := ssp.NewTokenGenerator(anyKey)

	t.Run("ReturnsAToken", func(t *testing.T) {
		got := generator.Token("someUser")
		t.Logf("Generated token: %s", got)
		assert.NotEmpty(t, got)
	})

	t.Run("Base64EncodesTheToken", func(t *testing.T) {
		got := generator.Token("someUser")
		_, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(got)
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
	generator := ssp.NewTokenGenerator(anyKey)
	currentTime, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", "2019-06-22 15:35:06 +0100 BST")
	generator.NowFunc = func() time.Time { return currentTime }

	t.Run("ReturnsNoErrorForValidToken", func(t *testing.T) {
		validToken := "hq84C7vMmeYtNPo5oEdsSDI2rCjpWAvs-U04DvVuQ79uVVBvHq--MwabFJVfbbI"
		err := generator.ValidateToken(validToken)
		assert.Nil(t, err)
	})

	t.Run("ReturnsTokenExpiredForOldToken", func(t *testing.T) {
		expiredToken := "Mq1C2bWBM7-gld_1bj_h7yZiL0OLggOxkuq6KXD7JTcuuykuZ0DljmpzTKpZBv8"
		err := generator.ValidateToken(expiredToken)
		assert.Equal(t, ssp.ErrTokenExpired, err)
	})

	t.Run("ReturnsTokenFormatInvalidForEmptyToken", func(t *testing.T) {
		err := generator.ValidateToken("")
		assert.Equal(t, ssp.ErrTokenFormatInvalid, err)
	})

	t.Run("ReturnsTokenFormatInvalidForRandomString", func(t *testing.T) {
		err := generator.ValidateToken("someinvalidtoken")
		assert.Equal(t, ssp.ErrTokenFormatInvalid, err)
	})
}
