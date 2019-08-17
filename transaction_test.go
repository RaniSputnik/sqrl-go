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
		Cmd: sqrl.CmdQuery,
		Idk: alice,
		Opt: []sqrl.Opt{},
	}
	validClient, _ := c.Encode()
	validServer := sqrl.Base64.EncodeToString([]byte("sqrl://example.com/sqrl?nut=123456789"))
	validIds := signature(aliceSig, validClient+validServer)

	newResponse := func() *sqrl.ServerMsg {
		return &sqrl.ServerMsg{ /* TODO */ }
	}

	t.Run("FailsWhenClientParamIsMissing", func(t *testing.T) {
		req := &sqrl.Request{
			Server:   validServer,
			Ids:      validIds,
			ClientIP: "10.0.0.1",
		}

		_, err := sqrl.Verify(req, nil, newResponse())
		assert.Equal(t, sqrl.ErrInvalidClient, err)
	})

	t.Run("FailsWhenServerParamIsMissing", func(t *testing.T) {
		req := &sqrl.Request{
			Client:   validClient,
			Ids:      validIds,
			ClientIP: "10.0.0.1",
		}

		_, err := sqrl.Verify(req, nil, newResponse())
		assert.Equal(t, sqrl.ErrInvalidServer, err)
	})

	// TODO: Do we fail when there's no previous transaction
	// and the cmd is ident? Shouldn't there always be a query first?

	t.Run("FailsWhenIDSInvalid", func(t *testing.T) {
		req := &sqrl.Request{
			Client:   validClient,
			Server:   validServer,
			Ids:      "invalid-sig",
			ClientIP: "10.0.0.1",
		}

		_, err := sqrl.Verify(req, nil, newResponse())
		assert.Equal(t, sqrl.ErrInvalidIDSig, err)
	})

	t.Run("FailsWhenSignedWrongPayload", func(t *testing.T) {
		wrongPayload := validClient + "-" + validServer
		invalidIds := signature(aliceSig, wrongPayload)

		req := &sqrl.Request{
			Client:   validClient,
			Server:   validServer,
			Ids:      invalidIds,
			ClientIP: "10.0.0.1",
		}

		_, err := sqrl.Verify(req, nil, newResponse())
		assert.Equal(t, sqrl.ErrInvalidIDSig, err)
	})

	t.Run("ReturnsParsedClientForAValidRequest", func(t *testing.T) {
		req := &sqrl.Request{
			Client:   validClient,
			Server:   validServer,
			Ids:      validIds,
			ClientIP: "10.0.0.1",
		}

		gotClient, err := sqrl.Verify(req, nil, newResponse())
		assert.NoError(t, err)
		assert.Equal(t, *c, *gotClient)
	})
}

func TestVerifyWithPreviousTransaction(t *testing.T) {
	alice, aliceSig := newIDKey()

	clientQuery := &sqrl.ClientMsg{
		Ver: []string{sqrl.V1},
		Cmd: sqrl.CmdQuery,
		Idk: alice,
		Opt: []sqrl.Opt{},
	}
	validClientQuery, _ := clientQuery.Encode()
	validServerQuery := sqrl.Base64.EncodeToString([]byte("sqrl://example.com/sqrl?nut=firstnut"))

	clientIdent := &sqrl.ClientMsg{
		Ver: []string{sqrl.V1},
		Cmd: sqrl.CmdIdent,
		Idk: alice,
		Opt: []sqrl.Opt{},
	}
	serverIdent := &sqrl.ServerMsg{
		Ver: []string{sqrl.V1},
		Nut: "secondnut",
		Tif: 0,
		Qry: "/sqrl?nut=secondnut",
	}
	validClientIdent, _ := clientIdent.Encode()
	validServerIdent, _ := serverIdent.Encode()

	newResponse := func() *sqrl.ServerMsg {
		return &sqrl.ServerMsg{ /* TODO */ }
	}

	t.Run("FailsWhenPreviousTransactionExistsButServerIsFullURL", func(t *testing.T) {
		prevRequest := &sqrl.Request{
			Client:   validClientQuery,
			Server:   validServerQuery,
			Ids:      signature(aliceSig, validClientQuery+validServerQuery),
			ClientIP: "10.0.0.1",
		}
		prevTransaction := &sqrl.Transaction{ /*TODO: Reply */ Request: prevRequest}

		reusedServerParam := validServerQuery
		req := &sqrl.Request{
			Client:   validClientIdent,
			Server:   reusedServerParam,
			Ids:      signature(aliceSig, validClientIdent+reusedServerParam),
			ClientIP: "10.0.0.1",
		}

		_, err := sqrl.Verify(req, prevTransaction, newResponse())
		assert.Equal(t, sqrl.ErrInvalidServer, err)
	})

	t.Run("FailsWhenClientIPDoesNotMatch", func(t *testing.T) {
		prevRequest := &sqrl.Request{
			Client:   validClientQuery,
			Server:   validServerQuery,
			Ids:      signature(aliceSig, validClientQuery+validServerQuery),
			ClientIP: "10.0.0.1",
		}
		prevTransaction := &sqrl.Transaction{ /*TODO: Reply */ Request: prevRequest}

		req := &sqrl.Request{
			Client:   validClientIdent,
			Server:   validServerIdent,
			Ids:      signature(aliceSig, validClientIdent+validServerIdent),
			ClientIP: "10.0.0.2",
		}

		_, err := sqrl.Verify(req, prevTransaction, newResponse())
		assert.Equal(t, sqrl.ErrIPMismatch, err)
	})

	t.Run("ReturnsNoErrorWhenClientIPDoesNotMatchButNoIPTestOptIsSet", func(t *testing.T) {
		clientQuery := &sqrl.ClientMsg{
			Ver: []string{sqrl.V1},
			Cmd: sqrl.CmdQuery,
			Idk: alice,
			Opt: []sqrl.Opt{sqrl.OptNoIPTest},
		}
		clientIdent := &sqrl.ClientMsg{
			Ver: []string{sqrl.V1},
			Cmd: sqrl.CmdIdent,
			Idk: alice,
			Opt: []sqrl.Opt{sqrl.OptNoIPTest},
		}
		validClientQuery, _ := clientQuery.Encode()
		validClientIdent, _ := clientIdent.Encode()

		prevRequest := &sqrl.Request{
			Client:   validClientQuery,
			Server:   validServerQuery,
			Ids:      signature(aliceSig, validClientQuery+validServerQuery),
			ClientIP: "10.0.0.1",
		}
		prevTransaction := &sqrl.Transaction{ /*TODO: Reply */ Request: prevRequest}

		req := &sqrl.Request{
			Client:   validClientIdent,
			Server:   validServerIdent,
			Ids:      signature(aliceSig, validClientIdent+validServerIdent),
			ClientIP: "10.0.0.2",
		}

		gotClient, err := sqrl.Verify(req, prevTransaction, newResponse())
		if assert.NoError(t, err) {
			assert.Equal(t, *clientIdent, *gotClient)
		}
	})

	t.Run("SetsTheIPMatchTIFlags", func(t *testing.T) {
		prevRequest := &sqrl.Request{
			Client:   validClientQuery,
			Server:   validServerQuery,
			Ids:      signature(aliceSig, validClientQuery+validServerQuery),
			ClientIP: "10.0.0.1",
		}
		prevTransaction := &sqrl.Transaction{ /*TODO: Reply */ Request: prevRequest}

		match := &sqrl.Request{
			Client:   validClientIdent,
			Server:   validServerIdent,
			Ids:      signature(aliceSig, validClientIdent+validServerIdent),
			ClientIP: "10.0.0.1",
		}

		nonMatch := &sqrl.Request{
			Client:   validClientIdent,
			Server:   validServerIdent,
			Ids:      signature(aliceSig, validClientIdent+validServerIdent),
			ClientIP: "10.0.0.2",
		}

		resMatch := newResponse()
		_, _ = sqrl.Verify(match, prevTransaction, resMatch)
		assert.True(t, resMatch.Is(sqrl.TIFIPMatch), "IP Match should be set")

		resNonMatch := newResponse()
		_, _ = sqrl.Verify(nonMatch, prevTransaction, resNonMatch)
		assert.False(t, resNonMatch.Is(sqrl.TIFIPMatch), "IP Match should not be set")
	})

	// TODO: Return an error when the client Opt have changed between requests

	// TODO: Return an error when the IDK have changed between requests

	t.Run("ReturnsParsedClientForAValidRequest", func(t *testing.T) {
		prevRequest := &sqrl.Request{
			Client:   validClientQuery,
			Server:   validServerQuery,
			Ids:      signature(aliceSig, validClientQuery+validServerQuery),
			ClientIP: "10.0.0.1",
		}
		prevTransaction := &sqrl.Transaction{ /*TODO: Reply */ Request: prevRequest}

		req := &sqrl.Request{
			Client:   validClientIdent,
			Server:   validServerIdent,
			Ids:      signature(aliceSig, validClientIdent+validServerIdent),
			ClientIP: "10.0.0.1",
		}

		gotClient, err := sqrl.Verify(req, prevTransaction, newResponse())
		if assert.NoError(t, err) {
			assert.Equal(t, *clientIdent, *gotClient)
		}
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
