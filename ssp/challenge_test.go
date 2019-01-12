package ssp_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/RaniSputnik/sqrl-go/ssp"
	"github.com/stretchr/testify/assert"
)

func TestGenerateChallenge(t *testing.T) {
	insecureKey := make([]byte, 16)
	mockServer := sqrl.Configure(insecureKey)
	getLoginRequest := httptest.NewRequest(http.MethodGet, "https://example.com/login", nil)

	t.Run("ReturnsAURLWithTheCorrectDomain", func(t *testing.T) {
		requestURL, _ := ssp.GenerateChallenge(mockServer, getLoginRequest)

		got, err := url.Parse(requestURL)
		assert.NoError(t, err)
		assert.Equal(t, "example.com", got.Host)
	})
}
