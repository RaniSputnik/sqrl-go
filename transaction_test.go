package sqrl_test

import (
	"testing"

	"github.com/RaniSputnik/sqrl-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
)

func TestVerifyFirstTransaction(t *testing.T) {
	alice, aliceSig := newIDKey()

	c := &sqrl.ClientMsg{
		Ver: []string{sqrl.V1},
		Cmd: sqrl.CmdIdent,
		Idk: alice,
		Opt: []sqrl.Opt{},
	}
	validClient, _ := c.Encode()
	validServer := sqrl.Base64.EncodeToString([]byte("sqrl://example.com/sqrl?nut=123456789"))
	validIds := signature(aliceSig, validClient+validServer)

	newResponse := func() *sqrl.ServerMsg {
		return &sqrl.ServerMsg{ /* TODO */ }
	}

	t.Run("FailsWhenIDSInvalid", func(t *testing.T) {
		transaction := &sqrl.Transaction{
			Client:   validClient,
			Server:   validServer,
			Ids:      "invalid-sig",
			ClientIP: "10.0.0.1",
		}

		_, err := sqrl.Verify(transaction, nil, newResponse())
		assert.Equal(t, sqrl.ErrInvalidIDSig, err)
	})

	t.Run("FailsWhenSignedWrongPayload", func(t *testing.T) {
		wrongPayload := validClient + "-" + validServer
		invalidIds := signature(aliceSig, wrongPayload)

		transaction := &sqrl.Transaction{
			Client:   validClient,
			Server:   validServer,
			Ids:      invalidIds,
			ClientIP: "10.0.0.1",
		}

		_, err := sqrl.Verify(transaction, nil, newResponse())
		assert.Equal(t, sqrl.ErrInvalidIDSig, err)
	})

	t.Run("ReturnsParsedClientForAValidRequest", func(t *testing.T) {
		transaction := &sqrl.Transaction{
			Client:   validClient,
			Server:   validServer,
			Ids:      validIds,
			ClientIP: "10.0.0.1",
		}

		gotClient, err := sqrl.Verify(transaction, nil, newResponse())
		assert.NoError(t, err)
		assert.Equal(t, *c, *gotClient)
	})
}

func newIDKey() (sqrl.Identity, []byte) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	return sqrl.Identity(sqrl.Base64.EncodeToString(pub)), priv
}

func signature(privateKey ed25519.PrivateKey, payload string) sqrl.Signature {
	sig := sqrl.Base64.EncodeToString(ed25519.Sign(privateKey, []byte(payload)))
	return sqrl.Signature(sig)
}
