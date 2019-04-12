package sqrl_test

import (
	"testing"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/stretchr/testify/assert"
)

func TestServerRedirectURL(t *testing.T) {
	anyKey := make([]byte, 16)

	t.Run("DefaultsToIndexURL", func(t *testing.T) {
		s := sqrl.Configure(anyKey)
		assert.Equal(t, s.RedirectURL(), "/")
	})

	t.Run("RespectsTheSpecifiedRedirectURL", func(t *testing.T) {
		s := sqrl.Configure(anyKey).WithRedirectURL("https://example.com")
		assert.Equal(t, s.RedirectURL(), "https://example.com")
	})

	t.Run("HandlesURLsWithPaths", func(t *testing.T) {
		s := sqrl.Configure(anyKey).WithRedirectURL("https://example.com/foo/bar")
		assert.Equal(t, s.RedirectURL(), "https://example.com/foo/bar")
	})

	t.Run("StripsAnyQueryStringParams", func(t *testing.T) {
		s := sqrl.Configure(anyKey).WithRedirectURL("https://example.com?foo=bar")
		assert.Equal(t, s.RedirectURL(), "https://example.com")
	})
}
