package sqrl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

func TestNut(t *testing.T) {
	anyKey := make([]byte, 16)
	nutter := sqrl.NewNutter(anyKey)

	t.Run("ReturnsANonEmptyValue", func(t *testing.T) {
		result := nutter.Next()
		assert.NotEmpty(t, result)
	})

	t.Run("ReturnsBase64EncodedString", func(t *testing.T) {
		result := nutter.Next()
		_, err := sqrl.Base64.DecodeString(string(result))
		assert.NoError(t, err, "Expected nut to be base64 encoded, but got an error during decoding")
	})

	t.Run("DoesNotRepeatNuts", func(t *testing.T) {
		results := map[sqrl.Nut]struct{}{}

		for i := 0; i < 100; i++ {
			result := nutter.Next()
			if !assert.NotContainsf(t, results, result, "Found duplicate nut: '%s'", result) {
				break
			}
			results[result] = struct{}{}
			t.Logf("Got nut: %s", result)
		}
	})

	// TODO: Check allows concurrent generation of nuts?
}
