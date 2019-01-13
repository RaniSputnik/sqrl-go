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
	fatal(t, assert.NoError(t, err, "Expected no HTTP/connection error"))
	defer res.Body.Close()

	// Assert headers
	assert.Equal(t, http.StatusOK, res.StatusCode, "Expected successful status code")
	assert.True(t, strings.HasPrefix("application/json", res.Header.Get("Content-Type")), "Expected response to have Content-Type 'application/json'")

	// Assert response body
	type nutRes struct {
		Nut string `json:"nut"`
	}
	var got nutRes
	err = json.NewDecoder(res.Body).Decode(&got)
	fatal(t, assert.NoError(t, err, "Expected body to be decoded as JSON successfully"))
	assert.NotEmpty(t, got.Nut, "Expected nut to be returned")
}

// TODO: Test the same nut is not returned twice

func fatal(t *testing.T, ok bool) {
	if !ok {
		t.FailNow()
	}
}
