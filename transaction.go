package sqrl

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrInvalidIDSig = errors.New("invalid identity signature")
)

type Transaction struct {
	Client string
	Server string
	Ids    Signature
	// TODO: Pids
	// TODO: Suk

	ClientIP string
}

// Verify checks that a transaction (request from a SQRL client) is valid
//
// A previous transaction should be provided if one exists. If no previous
// transaction is provided, the request is presumed to be the first transaction.
//
// If an error is encoutered, the precise error will be returned and the
// correct transaction information flags will be set on the response.
func Verify(t *Transaction, prev *Transaction, response *ServerMsg) error {
	client, errc := ParseClient(t.Client)
	if errc != nil {
		response.Tif = response.Tif | TIFCommandFailed | TIFClientFailure
		return errors.New("invalid client param")
	}
	serverOK := verifyServer(t.Server, prev)
	if !serverOK {
		response.Tif = response.Tif | TIFCommandFailed | TIFClientFailure
		return errors.New("invalid server param")
	}
	signedPayload := t.Client + t.Server
	if !t.Ids.Verify(client.Idk, signedPayload) {
		response.Tif = response.Tif | TIFCommandFailed | TIFClientFailure
		return ErrInvalidIDSig
	}

	// TODO: Verify IP Match

	// TODO: Verify IDK Match

	return nil
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
		msg, err := ParseServer(serverRaw)
		if err != nil || msg == nil || msg.Nut == "" {
			return false
		}
		return true
	}
}
