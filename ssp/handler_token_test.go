package ssp_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RaniSputnik/sqrl-go/ssp"
	"github.com/stretchr/testify/assert"
)

func TestTokenHandler(t *testing.T) {
	runHandler := func(h http.Handler) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodGet, "/token", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}

	t.Run("FailsWith401IfAuthFunctionRejects", func(t *testing.T) {
		rejectAll := func(_ *http.Request) error { return errors.New("reject everyone") }
		tokens := ssp.DefaultExchange(make([]byte, 16), time.Minute)

		h := anyServer().WithAuthentication(rejectAll).TokenHandler(tokens)
		result := runHandler(h)

		assert.Equal(t, http.StatusUnauthorized, result.Code)
	})

	t.Run("FailsWith404IfTokenInvalid", func(t *testing.T) {
		validator := NewTokenValidator().ReturnsValidationError(ssp.ErrTokenFormatInvalid)

		h := anyServer().TokenHandler(validator)
		result := runHandler(h)

		assert.Equal(t, http.StatusNotFound, result.Code)
	})

	t.Run("Returns200IfTokenValid", func(t *testing.T) {
		validator := NewTokenValidator().ReturnsUser("someUser")

		h := anyServer().TokenHandler(validator)
		result := runHandler(h)

		assert.Equal(t, http.StatusOK, result.Code)
	})

	t.Run("ReturnsUserFromTokenIfTokenValid", func(t *testing.T) {
		type tokenResponse struct {
			User string `json:"user"`
		}

		expectedUser := "someUser"
		validator := NewTokenValidator().ReturnsUser(expectedUser)

		h := anyServer().TokenHandler(validator)
		result := runHandler(h)

		var response tokenResponse
		err := json.NewDecoder(result.Body).Decode(&response)

		if assert.Nil(t, err) {
			assert.Equal(t, expectedUser, response.User)
		}
	})
}

type mockTokenValidator struct {
	Func struct {
		Validate struct {
			CalledWith struct {
				Token ssp.Token
			}
			Returns struct {
				UserID string
				Err    error
			}
		}
	}
}

func (m *mockTokenValidator) Validate(token ssp.Token) (string, error) {
	m.Func.Validate.CalledWith.Token = token
	return m.Func.Validate.Returns.UserID, m.Func.Validate.Returns.Err
}

func (m *mockTokenValidator) ReturnsValidationError(err error) *mockTokenValidator {
	m.Func.Validate.Returns.UserID = ""
	m.Func.Validate.Returns.Err = err
	return m
}

func (m *mockTokenValidator) ReturnsUser(userID string) *mockTokenValidator {
	m.Func.Validate.Returns.UserID = userID
	m.Func.Validate.Returns.Err = nil
	return m
}

func NewTokenValidator() *mockTokenValidator {
	return &mockTokenValidator{}
}
