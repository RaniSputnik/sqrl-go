package client_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/RaniSputnik/sqrl-go/client"
)

const serverResponseKnownUser = "dmVyPTENCm51dD01aHFaS3VIeXE1dDZ5Mmlmb1czd1B3DQp0aWY9NQ0KcXJ5PS9zcXJsP251dD01aHFaS3VIeXE1dDZ5Mmlmb1czd1B3DQo"

func TestLogin(t *testing.T) {
	t.Run("RejectsEmptyUrl", func(t *testing.T) {
		invalidUri := ":"
		expectErr(t, client.ErrUriInvalid, client.Login(invalidUri))
	})

	t.Run("RejectsUriWithoutSQRLProtocol", func(t *testing.T) {
		expectErr(t, client.ErrUriInvalid, client.Login("https://example.com"))
	})

	t.Run("PostsQueryRequestToServer", func(t *testing.T) {
		var receivedRequest *http.Request
		s := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedRequest = r
			//nolint:errcheck
			w.Write([]byte(serverResponseKnownUser))
		}))
		defer s.Close()
		client.HttpClient = s.Client()

		serverURL, _ := url.Parse(s.URL)
		serverURL.Scheme = "sqrl"
		sqrlUri := serverURL.String()

		t.Logf("Making request to SQRL server: '%s'", sqrlUri)
		expectErr(t, nil, client.Login(sqrlUri))

		if receivedRequest == nil {
			t.Errorf("Expected request to test server, but it was never made")
		}
	})
}

func expectErr(t *testing.T, expect, got error) {
	if got != expect {
		t.Errorf("Expected error: '%v', got: '%v'", expect, got)
	}
}
