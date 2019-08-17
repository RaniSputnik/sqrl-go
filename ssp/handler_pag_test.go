package ssp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/RaniSputnik/sqrl-go/ssp"
	"github.com/stretchr/testify/assert"
)

func TestPagHandler(t *testing.T) {
	callbackURL := "http://example.com/auth/callback"
	s := ssp.Configure(make([]byte, 16), callbackURL)

	validClientIP := "36.0.0.1"
	invalidClientIP := "36.0.0.2"
	validTransaction := &sqrl.Transaction{
		Request: &sqrl.Request{
			ClientIP: validClientIP,
		},
	}

	t.Run("ReturnsNotFoundWhenTheClientIPDoesNotMatch", func(t *testing.T) {
		mockStore := NewStore()
		mockStore.Func.GetFirstTransaction.Returns.Transaction = validTransaction
		mockStore.Func.GetIdentSuccess.Returns.Token = "sometoken"

		r := httptest.NewRequest(http.MethodGet, "/pag.sqrl?nut=123456789", nil)
		r.Header.Set("X-Forwarded-For", invalidClientIP)
		w := httptest.NewRecorder()
		h := s.PagHandler(mockStore)
		h.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("ReturnsTheRedirectURLWithToken", func(t *testing.T) {
		mockStore := NewStore()
		mockStore.Func.GetFirstTransaction.Returns.Transaction = validTransaction
		mockStore.Func.GetIdentSuccess.Returns.Token = "sometoken"

		r := httptest.NewRequest(http.MethodGet, "/pag.sqrl?nut=123456789", nil)
		r.Header.Set("X-Forwarded-For", validClientIP)
		w := httptest.NewRecorder()
		h := s.PagHandler(mockStore)
		h.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "http://example.com/auth/callback?token=sometoken", w.Body.String())
	})
}
