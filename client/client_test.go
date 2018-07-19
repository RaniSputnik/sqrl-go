package client_test

import (
	"testing"

	"github.com/RaniSputnik/sqrl-go/client"
)

func TestLogin(t *testing.T) {
	t.Run("RejectsEmptyUrl", func(t *testing.T) {
		invalidUri := ":"
		expectErr(t, client.ErrUriInvalid, client.Login(invalidUri))
	})

	t.Run("RejectsUriWithoutSQRLProtocol", func(t *testing.T) {
		expectErr(t, client.ErrUriInvalid, client.Login("https://example.com"))
	})
}

func expectErr(t *testing.T, expect, got error) {
	if got != expect {
		t.Errorf("Expected error: '%v', got: '%v'", expect, got)
	}
}
