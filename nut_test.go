package sqrl_test

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

func TestNext(t *testing.T) {
	anyRequest := httptest.NewRequest("GET", "/sqrl", nil)

	t.Run("ReturnsANonEmptyValue", func(t *testing.T) {
		result := sqrl.Next(anyRequest)
		assert.NotEmpty(t, result)
	})

	t.Run("ReturnsBase64EncodedString", func(t *testing.T) {
		result := sqrl.Next(anyRequest)
		_, err := sqrl.Base64.DecodeString(string(result))
		assert.NoError(t, err, "Expected nut to be base64 encoded, but got an error during decoding")
	})

	t.Run("DoesNotRepeatNuts", func(t *testing.T) {
		results := map[sqrl.Nut]struct{}{}

		for i := 0; i < 100; i++ {
			result := sqrl.Next(anyRequest)
			if !assert.NotContainsf(t, results, result, "Found duplicate nut: '%s'", result) {
				break
			}
			results[result] = struct{}{}
		}
	})
}
