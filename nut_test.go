package sqrl_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

func TestNut(t *testing.T) {
	anyRequest := httptest.NewRequest("GET", "/sqrl", nil)
	anyKey := make([]byte, 16)
	server := sqrl.Configure(anyKey)

	t.Run("ReturnsANonEmptyValue", func(t *testing.T) {
		result := server.Nut(anyRequest)
		assert.NotEmpty(t, result)
	})

	t.Run("ReturnsBase64EncodedString", func(t *testing.T) {
		result := server.Nut(anyRequest)
		_, err := sqrl.Base64.DecodeString(string(result))
		assert.NoError(t, err, "Expected nut to be base64 encoded, but got an error during decoding")
	})

	t.Run("DoesNotRepeatNuts", func(t *testing.T) {
		results := map[sqrl.Nut]struct{}{}

		for i := 0; i < 100; i++ {
			result := server.Nut(anyRequest)
			if !assert.NotContainsf(t, results, result, "Found duplicate nut: '%s'", result) {
				break
			}
			results[result] = struct{}{}
		}
	})
}

func TestValidate(t *testing.T) {
	firstRequest := httptest.NewRequest("GET", "/sqrl", nil)
	firstRequest.RemoteAddr = "12.34.56.7"
	validRequest := httptest.NewRequest("GET", "/sqrl", nil)
	validRequest.RemoteAddr = "12.34.56.7"
	requestFromIncorrectIP := httptest.NewRequest("GET", "/sqrl", nil)
	requestFromIncorrectIP.RemoteAddr = "76.54.32.1"
	validRequestWithNoIP := httptest.NewRequest("GET", "/sqrl", nil)
	validRequestWithNoIP.RemoteAddr = ""

	anyKey := make([]byte, 16)
	server := sqrl.Configure(anyKey)
	validNut := server.Nut(firstRequest)
	validNutWithNoIPCheck := server.Nut(validRequestWithNoIP)

	t.Run("ReturnsFalseWhenNutIsEmpty", func(t *testing.T) {
		assert.False(t, server.Validate(sqrl.Nut(""), validRequest))
	})

	t.Run("ReturnsFalseWhenNutIsInvalid", func(t *testing.T) {
		assert.False(t, server.Validate(sqrl.Nut("invalid"), validRequest))
	})

	t.Run("ReturnsFalseWhenIPsDoNotMatch", func(t *testing.T) {
		assert.False(t, server.Validate(validNut, requestFromIncorrectIP))
	})

	t.Run("IgnoresIPCheckWhenIPBytesAreNotSet", func(t *testing.T) {
		assert.True(t, server.Validate(validNutWithNoIPCheck, validRequest))
	})

	t.Run("ReturnsFalseWhenNutHasExpired", func(t *testing.T) {
		shortExpiry := time.Millisecond * 5
		serverWithShortExpiry := sqrl.Configure(anyKey).WithNutExpiry(shortExpiry)
		validNutWithShortExpiry := serverWithShortExpiry.Nut(firstRequest)

		// Wait for the nut to expire
		<-time.After(time.Millisecond * 10)

		assert.False(t, serverWithShortExpiry.Validate(validNutWithShortExpiry, validRequest))
	})

	t.Run("ReturnsTrueWhenIPsMatchAndNutIsValid", func(t *testing.T) {
		assert.True(t, server.Validate(validNut, validRequest))
	})
}
