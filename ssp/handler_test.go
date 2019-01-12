package ssp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RaniSputnik/sqrl-go/ssp"
)

func TestServerNut(t *testing.T) {
	s := httptest.NewServer(ssp.Handler())
	res, err := http.Get(s.URL + "/nut.sqrl")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode, "Expected successful status code")
}
