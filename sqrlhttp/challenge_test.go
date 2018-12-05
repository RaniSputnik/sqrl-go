package sqrlhttp_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/RaniSputnik/sqrl-go"
	"github.com/stretchr/testify/assert"

	"github.com/RaniSputnik/sqrl-go/sqrlhttp"
)

func TestGenerateChallenge(t *testing.T) {
	insecureKey := make([]byte, 16)
	mockServer := sqrl.Configure(insecureKey)
	getLoginRequest := httptest.NewRequest(http.MethodGet, "https://example.com/login", nil)

	t.Run("ReturnsAURLWithTheCorrectDomain", func(t *testing.T) {
		requestURL, _ := sqrlhttp.GenerateChallenge(mockServer, getLoginRequest)

		got, err := url.Parse(requestURL)
		assert.NoError(t, err)
		assert.Equal(t, "example.com", got.Host)
	})
}
