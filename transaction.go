package sqrl

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrInvalidClient = errors.New("invalid client param")
	ErrInvalidServer = errors.New("invalid server param")
	ErrInvalidIDSig  = errors.New("invalid identity signature")
	ErrIPMismatch    = errors.New("ip does not match")
)

type Transaction struct {
	Nut    Nut
	Client string
	Server string
	Ids    Signature
	// TODO: Pids
	// TODO: Suk

	ClientIP string

	// TODO: Should we add response here?
	// This way we'd be able to set the response
	// on the current transaction during verification
	// and compare the server parameter against the
	// response on the previous transaction
	//
	// Reply *ServerMsg
}

// Verify checks that a transaction (request from a SQRL client) is valid
//
// A previous transaction should be provided if one exists. If no previous
// transaction is provided, the request is presumed to be the first transaction.
//
// If an error is encoutered, the precise error will be returned and the
// correct transaction information flags will be set on the response.
func Verify(t *Transaction, prev *Transaction, response *ServerMsg) (*ClientMsg, error) {
	client, errc := ParseClient(t.Client)
	if errc != nil {
		response.Tif = response.Tif | TIFCommandFailed | TIFClientFailure
		return nil, ErrInvalidClient
	}
	serverOK := verifyServer(t.Server, prev)
	if !serverOK {
		response.Tif = response.Tif | TIFCommandFailed | TIFClientFailure
		return nil, ErrInvalidServer
	}
	signedPayload := t.Client + t.Server
	if !t.Ids.Verify(client.Idk, signedPayload) {
		response.Tif = response.Tif | TIFCommandFailed | TIFClientFailure
		return nil, ErrInvalidIDSig
	}

	if prev == nil {
		return client, nil
	}

	ipMustMatch := !client.HasOpt(OptNoIPTest)
	if ipMustMatch && prev.ClientIP != t.ClientIP {
		return nil, ErrIPMismatch
	}

	// TODO: Verify IDK Match

	return client, nil
}

func verifyServer(serverRaw string, prev *Transaction) bool {
	// TODO: Here we accept EITHER a URL or ServerMsg
	// However we know that ONLY the first request
	// from the client should be a URL.
	// Is there a way for us to ensure that here?

	bytes, err := Base64.DecodeString(serverRaw)
	if err != nil {
		return false
	}

	server := string(bytes)
	if strings.HasPrefix(server, "sqrl") {
		if prev != nil {
			// Providing the previous query URL as the server
			// param is ONLY valid for the first transaction
			return false
		}

		serverURL, err := url.Parse(server)
		if err != nil {
			return false
		}

		// TODO: Assert URL matches server configuration
		// eg. domain, "server friendly name", etc.

		nut := serverURL.Query().Get("nut")
		if nut == "" {
			return false
		}
		return true

	} else {
		// TODO: We must be able to do something smarter with this
		// We should be able to check the previous server val we sent
		// to the client and compare this against that.
		msg, err := ParseServer(serverRaw)
		if err != nil || msg == nil || msg.Nut == "" {
			return false
		}
		return true
	}
}
