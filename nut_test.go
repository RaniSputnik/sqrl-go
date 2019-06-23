package sqrl_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

func TestNut(t *testing.T) {
	anyKey := make([]byte, 16)
	nutter := sqrl.NewNutter(anyKey)

	t.Run("ReturnsANonEmptyValue", func(t *testing.T) {
		result := nutter.Nut(sqrl.NoClientID)
		assert.NotEmpty(t, result)
	})

	t.Run("ReturnsBase64EncodedString", func(t *testing.T) {
		result := nutter.Nut(sqrl.NoClientID)
		_, err := sqrl.Base64.DecodeString(string(result))
		assert.NoError(t, err, "Expected nut to be base64 encoded, but got an error during decoding")
	})

	t.Run("DoesNotRepeatNuts", func(t *testing.T) {
		results := map[sqrl.Nut]struct{}{}

		for i := 0; i < 100; i++ {
			result := nutter.Nut(sqrl.NoClientID)
			if !assert.NotContainsf(t, results, result, "Found duplicate nut: '%s'", result) {
				break
			}
			results[result] = struct{}{}
		}
	})

	// TODO: Check allows concurrent generation of nuts?
}

func TestValidate(t *testing.T) {
	validIP := "12.34.56.7"
	anyKey := make([]byte, 16)
	nutter := sqrl.NewNutter(anyKey)
	validNut := nutter.Nut(validIP)

	t.Run("ReturnsFalseWhenNutIsEmpty", func(t *testing.T) {
		assert.False(t, nutter.Validate(sqrl.Nut(""), validIP))
	})

	t.Run("ReturnsFalseWhenNutIsInvalid", func(t *testing.T) {
		assert.False(t, nutter.Validate(sqrl.Nut("invalid"), validIP))
	})

	t.Run("ReturnsFalseWhenIPsDoNotMatch", func(t *testing.T) {
		incorrectIP := "76.54.32.1"
		assert.False(t, nutter.Validate(validNut, incorrectIP))
	})

	t.Run("IgnoresIPCheckWhenIPBytesAreNotSet", func(t *testing.T) {
		validNutWithNoIPCheck := nutter.Nut(sqrl.NoClientID)
		assert.True(t, nutter.Validate(validNutWithNoIPCheck, validIP))
	})

	t.Run("ReturnsFalseWhenComplexClientIDDoesNotMatch", func(t *testing.T) {
		complexClientID := "12.34.56.7+Chrome@70+uid:1234567"
		nonmatchClientID := "12.34.56.7+Safari+uid:1234567"
		nutWithComplexClient := nutter.Nut(complexClientID)
		assert.False(t, nutter.Validate(nutWithComplexClient, nonmatchClientID))
	})

	t.Run("ReturnsFalseWhenNoClientIDDoesNotMatch", func(t *testing.T) {
		// This case would defeat the client ID check if not implemented correctly
		// We should not be able to provide 'NoClientID' when validating a nut
		// that was created with a valid client id.
		assert.False(t, nutter.Validate(validNut, sqrl.NoClientID))
	})

	t.Run("ReturnsFalseWhenNutHasExpired", func(t *testing.T) {
		shortExpiry := time.Millisecond * 5
		nutterWithShortExpiry := sqrl.NewNutter(anyKey)
		nutterWithShortExpiry.Expiry = shortExpiry
		validNutWithShortExpiry := nutterWithShortExpiry.Nut(validIP)

		// Wait for the nut to expire
		<-time.After(time.Millisecond * 10)

		assert.False(t, nutterWithShortExpiry.Validate(validNutWithShortExpiry, validIP))
	})

	t.Run("ReturnsTrueWhenIPsMatchAndNutIsValid", func(t *testing.T) {
		assert.True(t, nutter.Validate(validNut, validIP))
	})
}
