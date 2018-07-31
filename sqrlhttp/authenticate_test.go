package sqrlhttp_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RaniSputnik/sqrl-go/sqrlhttp"
)

func TestAuthenticate(t *testing.T) {
	const emptyBody = ""
	const validServer = "whatever"

	t.Run("ReturnsBadRequestWhenContentTypeIsNotFormEncoded", func(t *testing.T) {
		h := sqrlhttp.Authenticate()
		w, r := setupAuthenticate(emptyBody)
		r.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(w, r)

		if expected := http.StatusBadRequest; w.Code != expected {
			t.Errorf("Expected: %d, got: %d", expected, w.Code)
		}
	})

	t.Run("ReturnsBadRequestWhenServerParamIsMissing", func(t *testing.T) {

	})

	t.Run("ReturnsBadRequestWhenServerParamIsInvalid", func(t *testing.T) {

	})

	t.Run("ReturnsBadRequestWhenClientParamIsMissing", func(t *testing.T) {
		h := sqrlhttp.Authenticate()
		w, r := setupAuthenticate(fmt.Sprintf("server=%s", validServer))
		h.ServeHTTP(w, r)

		if expected := http.StatusBadRequest; w.Code != expected {
			t.Errorf("Expected: %d, got: %d", expected, w.Code)
		}
	})

	t.Run("ReturnsBadRequestWhenClientStringIsInvalid", func(t *testing.T) {
		cases := []struct {
			Name  string
			Input string
		}{
			{"Empty", ""},
			{"VersionOnly", b64("ver=1\n\n\n\n")},
			{"Rubbish", b64("this is rubbish")},
			{"DuplicateVer", b64("ver=1\nver=2\nver=3\ncmd=query\n")},
			{"VerComesSecond", b64("cmd=query\nver=1")},
		}

		h := sqrlhttp.Authenticate()

		const expected = http.StatusBadRequest

		for _, test := range cases {
			t.Run(test.Name, func(t *testing.T) {
				w, r := setupAuthenticate(fmt.Sprintf("server=%s&client=%s", validServer, test.Input))
				h.ServeHTTP(w, r)
				if w.Code != expected {
					t.Errorf("Expected: %d, Got: %d", expected, w.Code)
				}
			})
		}
	})
}
func b64(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}

func setupAuthenticate(body string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set("User-Agent", "SQRL/1")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return w, r
}
