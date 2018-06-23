package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	sqlhttp "github.com/RaniSputnik/sqrl-go/http"
)

func TestAuthenticate(t *testing.T) {
	t.Run("ReturnsBadRequestWhenUserAgentNotProvided", func(t *testing.T) {
		h := sqlhttp.Authenticate()
		w, r := setupAuthenticate()
		r.Header.Del("User-Agent")

		h.ServeHTTP(w, r)

		if expected := http.StatusBadRequest; w.Code != expected {
			t.Errorf("Expected: %d, got: %d", expected, w.Code)
		}
	})

	t.Run("ReturnsBadRequestWhenUserAgentIsNotSQRL1", func(t *testing.T) {
		h := sqlhttp.Authenticate()
		w, r := setupAuthenticate()
		r.Header.Set("User-Agent", "SQRL/6")

		h.ServeHTTP(w, r)

		if expected := http.StatusBadRequest; w.Code != expected {
			t.Errorf("Expected: %d, got: %d", expected, w.Code)
		}
	})
}

func setupAuthenticate() (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("User-Agent", "SQRL/1")
	return w, r
}
