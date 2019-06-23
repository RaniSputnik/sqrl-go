package ssp_test

import (
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RaniSputnik/sqrl-go/ssp"
)

const validNut = "rNRqu8olcWLAPaDvsL4b6owTVfryjzbre3hWHWnNTrK_hIS_KgIDFt2eBDc"

func TestHandlerNutIsReturned(t *testing.T) {
	s := httptest.NewServer(ssp.Handler(anyServer(), noProtection()))
	res, err := http.Get(s.URL + "/nut.json")

	// Assert no errors
	fatal(t, assert.NoError(t, err,
		"Expected no HTTP/connection error"))
	defer res.Body.Close()

	// Assert headers
	assert.Equal(t, http.StatusOK, res.StatusCode,
		"Expected successful status code")
	assert.True(t, strings.HasPrefix("application/json", res.Header.Get("Content-Type")),
		"Expected response to have Content-Type 'application/json'")

	// Assert response body
	type nutRes struct {
		Nut string `json:"nut"`
	}
	var got nutRes
	err = json.NewDecoder(res.Body).Decode(&got)
	fatal(t, assert.NoError(t, err,
		"Expected body to be decoded as JSON successfully"))
	assert.NotEmpty(t, got.Nut,
		"Expected nut to be returned")
}

func TestHandlerNutIsUnique(t *testing.T) {
	s := httptest.NewServer(ssp.Handler(anyServer(), noProtection()))
	endpoint := s.URL + "/nut.json"

	type nutRes struct {
		Nut string `json:"nut"`
	}

	results := make(map[string]bool)

	for i := 0; i < 10; i++ {
		res, err := http.Get(endpoint)
		fatal(t, assert.NoError(t, err,
			"Expected no HTTP/connection error"))
		defer res.Body.Close()

		var got nutRes
		fatal(t, assert.NoError(t, json.NewDecoder(res.Body).Decode(&got),
			"Expected body to be decoded as JSON successfully"))

		seenBefore := results[got.Nut]
		fatal(t, assert.False(t, seenBefore, "Duplicate nut returned: '%s'", got.Nut))
		results[got.Nut] = true
	}
}

func TestQRCodeIsReturned(t *testing.T) {
	s := httptest.NewServer(ssp.Handler(anyServer(), noProtection()))
	res, err := http.Get(s.URL + "/qr.png?nut=" + validNut)

	// Assert no errors
	fatal(t, assert.NoError(t, err,
		"Expected no HTTP/connection error"))
	defer res.Body.Close()

	// Assert headers
	assert.Equal(t, http.StatusOK, res.StatusCode,
		"Expected successful status code")
	assert.Equal(t, "image/png", res.Header.Get("Content-Type"),
		"Expected response to have Content-Type 'image/png'")

	// Assert response body

	_, err = png.Decode(res.Body)
	fatal(t, assert.NoError(t, err,
		"Expected to decode the body as a PNG image successfully"))

	// TODO: We should compare images here to ensure the data was encoded successfully
}

func TestQRCodeIsReturnedAtSpecifiedSize(t *testing.T) {
	const givenSize = 64

	s := httptest.NewServer(ssp.Handler(anyServer(), noProtection()))
	res, err := http.Get(fmt.Sprintf("%s/qr.png?nut=%s&size=%d", s.URL, validNut, givenSize))
	fatal(t, assert.NoError(t, err,
		"Expected no HTTP/connection error"))
	defer res.Body.Close()

	img, err := png.Decode(res.Body)
	fatal(t, assert.NoError(t, err,
		"Expected to decode the body as a PNG image successfully"))

	size := img.Bounds().Size()
	assert.Equal(t, givenSize, size.X, "Expected image width to match 'size' query parameter")
	assert.Equal(t, givenSize, size.Y, "Expected image width to match 'size' query parameter")
}

func fatal(t *testing.T, ok bool) {
	if !ok {
		t.FailNow()
	}
}

func anyServer() *ssp.Server {
	return ssp.Configure(make([]byte, 16), "http://example.com/auth/callback")
}

func anyTokenGenerator() *ssp.TokenGenerator {
	return ssp.NewTokenGenerator(make([]byte, 16))
}

func noProtection() ssp.ServerToServerAuthValidationFunc {
	return func(r *http.Request) error {
		// Allow all server-to-server requests through
		// without authentication. Note: this should NEVER
		// be done in the wild, it's okay for testing.
		return nil
	}
}
