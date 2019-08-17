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

type Request struct {
	Nut    Nut
	Client string
	Server string
	Ids    Signature
	// TODO: Pids
	// TODO: Suk

	ClientIP string
}

type Transaction struct {
	Next Nut
	*Request
}

// Verify checks that a request from a SQRL client is valid.
//
// The transaction that started this session should be provided if one exists.
// If no previous transaction is provided, the request is presumed to be the
// first request for this session.
//
// Note: No attempt is made to verify the previous transaction (other than
// to compare it's properties to those of the new transaction). It is assumed
// that the previous transaction has already had it's signatures checked and
// payload validated.
//
// If a validation error is encoutered, the precise error will be returned and the
// correct transaction information flags will be set on the response.
func Verify(req *Request, first *Transaction, response *ServerMsg) (*ClientMsg, error) {
	if req.ClientIP == "" {
		// ClientIP MUST always be set correctly for same-device protections
		// to work correctly. We do not return an exported error here, because
		// this is not an error that should ever be caught in practise. It is
		// a development mistake and we check for it here to catch accidental
		// misuse that may result in a security vulnerability.
		return nil, errors.New("client ip should never be empty")
	}

	client, errc := ParseClient(req.Client)
	if errc != nil {
		response.Tif = response.Tif | TIFCommandFailed | TIFClientFailure
		return nil, ErrInvalidClient
	}
	serverOK := verifyServer(req.Server, first)
	if !serverOK {
		response.Tif = response.Tif | TIFCommandFailed | TIFClientFailure
		return nil, ErrInvalidServer
	}
	signedPayload := req.Client + req.Server
	if !req.Ids.Verify(client.Idk, signedPayload) {
		response.Tif = response.Tif | TIFCommandFailed | TIFClientFailure
		return nil, ErrInvalidIDSig
	}

	if first == nil {
		return client, nil
	}

	// TODO: Do we set IP Match for the first request? Presume not
	if first.ClientIP == req.ClientIP {
		response.Set(TIFIPMatch)
	}
	ipMustMatch := !client.HasOpt(OptNoIPTest)
	ipsMatch := response.Is(TIFIPMatch)
	if ipMustMatch && !ipsMatch {
		return nil, ErrIPMismatch
	}

	// TODO: Verify IDK Match

	// TODO: Is cmd "ident" allowed if there is no previous transaction?

	return client, nil
}

func verifyServer(serverRaw string, first *Transaction) bool {
	bytes, err := Base64.DecodeString(serverRaw)
	if err != nil {
		return false
	}

	server := string(bytes)
	if strings.HasPrefix(server, "sqrl") {
		if first != nil {
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
		return nut != ""

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
