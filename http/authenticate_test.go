package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sqlhttp "github.com/RaniSputnik/sqrl-go/http"
)

func TestAuthenticate(t *testing.T) {
	const emptyBody = ""

	t.Run("ReturnsBadRequestWhenUserAgentNotProvided", func(t *testing.T) {
		h := sqlhttp.Authenticate()
		w, r := setupAuthenticate(emptyBody)
		r.Header.Del("User-Agent")

		h.ServeHTTP(w, r)

		if expected := http.StatusBadRequest; w.Code != expected {
			t.Errorf("Expected: %d, got: %d", expected, w.Code)
		}
	})

	t.Run("ReturnsBadRequestWhenUserAgentIsNotSQRL1", func(t *testing.T) {
		h := sqlhttp.Authenticate()
		w, r := setupAuthenticate(emptyBody)
		r.Header.Set("User-Agent", "SQRL/6")

		h.ServeHTTP(w, r)

		if expected := http.StatusBadRequest; w.Code != expected {
			t.Errorf("Expected: %d, got: %d", expected, w.Code)
		}
	})

	t.Run("ReturnsBadRequestWhenContentTypeIsNotFormEncoded", func(t *testing.T) {
		h := sqlhttp.Authenticate()
		w, r := setupAuthenticate(emptyBody)
		r.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(w, r)

		if expected := http.StatusBadRequest; w.Code != expected {
			t.Errorf("Expected: %d, got: %d", expected, w.Code)
		}
	})
}

func setupAuthenticate(body string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set("User-Agent", "SQRL/1")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return w, r
}
