package ssp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RaniSputnik/sqrl-go/ssp"
)

var anyKey = make([]byte, 16)

func TestHandlerNutIsReturned(t *testing.T) {
	s := httptest.NewServer(ssp.Handler(anyKey))
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
	s := httptest.NewServer(ssp.Handler(anyKey))
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

func fatal(t *testing.T, ok bool) {
	if !ok {
		t.FailNow()
	}
}
