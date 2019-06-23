package ssp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RaniSputnik/sqrl-go/ssp"
	"github.com/stretchr/testify/assert"
)

func TestTokenHandler(t *testing.T) {
	t.Run("FailsWith401IfAuthFunctionRejects", func(t *testing.T) {
		rejectAll := func(_ *http.Request) error { return errors.New("reject everyone") }
		tokens := ssp.DefaultExchange(make([]byte, 16), time.Minute)

		h := anyServer().WithAuthentication(rejectAll).TokenHandler(tokens)
		r := httptest.NewRequest(http.MethodGet, "/token", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// TODO: FailsWith404IfTokenInvalid

	// TODO: Returns200IfTokenValid

	// TODO: ReturnsUserFromTokenIfTokenValid
}
