package sqrlhttp_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RaniSputnik/sqrl-go"

	"github.com/RaniSputnik/sqrl-go/sqrlhttp"
)

func TestAuthenticate(t *testing.T) {
	const emptyBody = ""
	const validServer = "dmVyPTENCm51dD11c2pxZmgzdFJoYWdHbjkyN0RZRmpRDQp0aWY9NA0KcXJ5PS9zcXJsP251dD11c2pxZmgzdFJoYWdHbjkyN0RZRmpRDQpzaW49MA0K"

	t.Run("ReturnsClientErrorWhenContentTypeIsNotFormEncoded", func(t *testing.T) {
		h := sqrlhttp.Authenticate()
		w, r := setupAuthenticate(emptyBody)
		r.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(w, r)

		got, err := sqrl.ParseServer(w.Body.String())
		if assert.NoError(t, err) {
			assert.True(t, got.Tif.Has(sqrl.TIFCommandFailed))
			assert.True(t, got.Tif.Has(sqrl.TIFClientFailure))
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

		got, err := sqrl.ParseServer(w.Body.String())
		if assert.NoError(t, err) {
			assert.True(t, got.Tif.Has(sqrl.TIFCommandFailed))
			assert.True(t, got.Tif.Has(sqrl.TIFClientFailure))
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

		for _, test := range cases {
			t.Run(test.Name, func(t *testing.T) {
				w, r := setupAuthenticate(fmt.Sprintf("server=%s&client=%s", validServer, test.Input))
				h.ServeHTTP(w, r)

				got, err := sqrl.ParseServer(w.Body.String())
				if assert.NoError(t, err) {
					assert.True(t, got.Tif.Has(sqrl.TIFCommandFailed))
					assert.True(t, got.Tif.Has(sqrl.TIFClientFailure))
				}
			})
		}
	})
}

func b64(in string) string {
	return sqrl.Base64.EncodeToString([]byte(in))
}

func setupAuthenticate(body string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set("User-Agent", "SQRL/1")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return w, r
}
