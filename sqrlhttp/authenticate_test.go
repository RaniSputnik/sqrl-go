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

const emptyBody = ""
const validServer = "dmVyPTENCm51dD11c2pxZmgzdFJoYWdHbjkyN0RZRmpRDQp0aWY9NA0KcXJ5PS9zcXJsP251dD11c2pxZmgzdFJoYWdHbjkyN0RZRmpRDQpzaW49MA0K"

// Sourced from grc.com SQRL diagnostic site
// https://www.grc.com/sqrl/logsample.htm
const validQueryBody = "client=dmVyPTENCmNtZD1xdWVyeQ0KaWRrPVpIa2RQTDM0eWFhSmR5aUtVT1F1SS1zMmtqei1uSGcwVU5RMFpBcjZlZHMNCg&server=c3FybDovL3d3dy5ncmMuY29tL3Nxcmw_bnV0PUNYam9xNVJla3FTNUQ1d3V5QktMUlEmc2ZuPVIxSkQ&ids=JqY1dMvWFunVSykecky3pM21KtW67gegPxcEpiA2obUzb1igxrLrEj5hI9QPZb8dIAnn8TtYSpPj4mRFFqNcAA"

func TestAuthenticateReturnsClientErrorWhenContentTypeIsNotFormEncoded(t *testing.T) {
	h := sqrlhttp.Authenticate(NewDelegate())
	w, r := setupAuthenticate(emptyBody)
	r.Header.Set("Content-Type", "application/json")

	h.ServeHTTP(w, r)

	got, err := sqrl.ParseServer(w.Body.String())
	if assert.NoError(t, err) {
		assert.True(t, got.Is(sqrl.TIFCommandFailed))
		assert.True(t, got.Is(sqrl.TIFClientFailure))
	}
}

// TODO: Invalid Server Param

// TODO: Invalid Server Param

func TestAuthenticateReturnsClientFailureWhenClientParamIsMissing(t *testing.T) {
	h := sqrlhttp.Authenticate(NewDelegate())
	w, r := setupAuthenticate(fmt.Sprintf("server=%s", validServer))
	h.ServeHTTP(w, r)

	got, err := sqrl.ParseServer(w.Body.String())
	if assert.NoError(t, err) {
		assert.True(t, got.Is(sqrl.TIFCommandFailed))
		assert.True(t, got.Is(sqrl.TIFClientFailure))
	}
}

func TestAuthenticateReturnsClientFailureWhenClientStringIsInvalid(t *testing.T) {
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

	h := sqrlhttp.Authenticate(NewDelegate())

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			w, r := setupAuthenticate(fmt.Sprintf("server=%s&client=%s", validServer, test.Input))
			h.ServeHTTP(w, r)

			got, err := sqrl.ParseServer(w.Body.String())
			if assert.NoError(t, err) {
				assert.True(t, got.Is(sqrl.TIFCommandFailed))
				assert.True(t, got.Is(sqrl.TIFClientFailure))
			}
		})
	}
}

func TestAuthenticateReturnsCurrentIDMatchWhenIDIsKnown(t *testing.T) {
	w, r := setupAuthenticate(validQueryBody)
	h := sqrlhttp.Authenticate(NewDelegate().ReturnsKnownIdentity())

	h.ServeHTTP(w, r)

	got, err := sqrl.ParseServer(w.Body.String())
	if assert.NoError(t, err) {
		assert.True(t, got.Is(sqrl.TIFCurrentIDMatch))
	}
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
