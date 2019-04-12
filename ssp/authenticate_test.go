package ssp_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/RaniSputnik/sqrl-go/ssp"
)

const emptyBody = ""
const validServer = "dmVyPTENCm51dD11c2pxZmgzdFJoYWdHbjkyN0RZRmpRDQp0aWY9NA0KcXJ5PS9zcXJsP251dD11c2pxZmgzdFJoYWdHbjkyN0RZRmpRDQpzaW49MA0K"

// Sourced from grc.com SQRL diagnostic site
// https://www.grc.com/sqrl/logsample.htm
const invalidQuerySignature = "client=dmVyPTENCmNtZD1xdWVyeQ0KaWRrPVpIa2RQTDM0eWFhSmR5aUtVT1F1SS1zMmtqei1uSGcwVU5RMFpBcjZlZHMNCg&server=c3FybDovL3d3dy5ncmMuY29tL3Nxcmw_bnV0PUNYam9xNVJla3FTNUQ1d3V5QktMUlEmc2ZuPVIxSkQ&ids=invalid"

const validQueryBody = "client=dmVyPTENCmNtZD1xdWVyeQ0KaWRrPVpIa2RQTDM0eWFhSmR5aUtVT1F1SS1zMmtqei1uSGcwVU5RMFpBcjZlZHMNCg&server=c3FybDovL3d3dy5ncmMuY29tL3Nxcmw_bnV0PUNYam9xNVJla3FTNUQ1d3V5QktMUlEmc2ZuPVIxSkQ&ids=JqY1dMvWFunVSykecky3pM21KtW67gegPxcEpiA2obUzb1igxrLrEj5hI9QPZb8dIAnn8TtYSpPj4mRFFqNcAA"

const validIdentBody = "client=dmVyPTENCmNtZD1pZGVudA0KaWRrPVpIa2RQTDM0eWFhSmR5aUtVT1F1SS1zMmtqei1uSGcwVU5RMFpBcjZlZHMNCg&server=dmVyPTENCm51dD01aHFaS3VIeXE1dDZ5Mmlmb1czd1B3DQp0aWY9NQ0KcXJ5PS9zcXJsP251dD01aHFaS3VIeXE1dDZ5Mmlmb1czd1B3DQo&ids=z__MvVTGpeDLLPj3O9QLNrkcvsk_8iuipu-DWalCfQWuP1xXom3HW1MhXNOYYhYiO2Kx2qMgT3D0uze3hdYLDg"

func TestAuthenticateReturnsClientErrorWhenContentTypeIsNotFormEncoded(t *testing.T) {
	h := ssp.Authenticate(anyServer(), NewDelegate())
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
	h := ssp.Authenticate(anyServer(), NewDelegate())
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

	h := ssp.Authenticate(anyServer(), NewDelegate())

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
	h := ssp.Authenticate(anyServer(), NewDelegate().ReturnsKnownIdentity())

	h.ServeHTTP(w, r)

	got, err := sqrl.ParseServer(w.Body.String())
	if assert.NoError(t, err) {
		assert.True(t, got.Is(sqrl.TIFCurrentIDMatch), "Expected 'TIFCurrentIDMatch' flag to be set")
	}
}

func TestAuthenticateReturnsNoIDMatchWhenIDIsUnknown(t *testing.T) {
	w, r := setupAuthenticate(validQueryBody)
	h := ssp.Authenticate(anyServer(), NewDelegate().ReturnsUnknownIdentity())

	h.ServeHTTP(w, r)

	got, err := sqrl.ParseServer(w.Body.String())
	if assert.NoError(t, err) {
		assert.False(t, got.Is(sqrl.TIFCurrentIDMatch), "Expected 'TIFCurrentIDMatch' flag to not be set")
	}
}

func TestAuthenticateReturnsClientErrorWhenSignatureInvalid(t *testing.T) {
	w, r := setupAuthenticate(invalidQuerySignature)
	h := ssp.Authenticate(anyServer(), NewDelegate())

	h.ServeHTTP(w, r)

	got, err := sqrl.ParseServer(w.Body.String())
	if assert.NoError(t, err) {
		assert.True(t, got.Is(sqrl.TIFClientFailure), "Expected client failure")
	}
}

func TestAuthenticateCallsDelegateVerifiedWhenIdentSuccessful(t *testing.T) {
	delegate := NewDelegate().ReturnsKnownIdentity()
	w, r := setupAuthenticate(validIdentBody)
	h := ssp.Authenticate(anyServer(), delegate)

	h.ServeHTTP(w, r)

	// TODO: this assertion fails in a really unintuitive way when
	// this test fails - see if we can improve this.
	assert.Equal(t,
		sqrl.Identity("ZHkdPL34yaaJdyiKUOQuI-s2kjz-nHg0UNQ0ZAr6eds"),
		delegate.Func.Verified.CalledWith.Id)
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

func anyServer() *sqrl.Server {
	return sqrl.Configure(make([]byte, 16))
}
