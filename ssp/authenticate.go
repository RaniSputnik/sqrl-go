package ssp

import (
	"fmt"
	"log"
	"net/http"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

const xFormURLEncoded = "application/x-www-form-urlencoded"

var v1Only = []string{sqrl.V1}

func clientFailure(response *sqrl.ServerMsg) {
	response.Tif = response.Tif | sqrl.TIFCommandFailed | sqrl.TIFClientFailure
}

func serverError(response *sqrl.ServerMsg) {
	response.Tif |= sqrl.TIFCommandFailed
}

func Authenticate(server *sqrl.Server, delegate Delegate) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got SQRL request: %v\n", r)

		// Reference implementation here
		// https://github.com/Novators/libsqrl/blob/c/src/server_protocol.c

		response := genNextResponse(server, r)
		defer writeResponse(w, response)

		if r.Header.Get("Content-Type") != xFormURLEncoded {
			clientFailure(response)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("Failed to parse form: %s\n", err)
			clientFailure(response)
			return
		}

		clientRaw := r.Form.Get("client")
		serverRaw := r.Form.Get("server")
		ids := sqrl.Signature(r.Form.Get("ids"))
		// TODO: pids

		client, errc := sqrl.ParseClient(clientRaw)
		// Ensure that it is the URL / params we sent
		//_, errs := sqrl.ParseServer(r.Form.Get("server"))
		if errc != nil /*|| errs != nil */ {
			clientFailure(response)
			return
		}

		// TODO: Verify server parameter in some way

		signedPayload := clientRaw + serverRaw
		if !ids.Verify(client.Idk, signedPayload) {
			clientFailure(response)
			return
		}
		// TODO: Verify previous identity signatures

		// TODO: Fetch user details
		// TODO: Fetch user from previous identity

		// TODO: Test for IP Match

		// TODO: Pass previous identities to "known"
		isKnown, err := delegate.Known(r.Context(), client.Idk)
		if err != nil {
			log.Printf("Failed to process response: %v\n", err)
			serverError(response)
			return
		} else if isKnown {
			response.Tif |= sqrl.TIFCurrentIDMatch
		}

		switch client.Cmd {
		case sqrl.CmdIdent:
			err := delegate.Authenticated(r.Context(), client.Idk)
			if err != nil {
				log.Fatalf("Failed to check authenticated: %v\n", err)
				serverError(response)
			}

			if client.HasOpt(sqrl.OptCPS) {
				// TODO: Configure the redirect URL for authentication
				token := "todo-token"
				response.URL = fmt.Sprintf("%s?%s", server.RedirectURL(), token)
			}

		case sqrl.CmdQuery:
			// Do nothing

		default:
			// In all other cases, not supported
			response.Tif |= sqrl.TIFFunctionNotSupported
		}
	})
}

func genNextResponse(server *sqrl.Server, r *http.Request) *sqrl.ServerMsg {
	// TODO: How to configure this endpoint?
	endpoint := "/sqrl" + r.URL.Path

	nextNut := server.Nut(clientID(r))
	return &sqrl.ServerMsg{
		Ver: v1Only,
		Nut: nextNut,
		Qry: endpoint + "?nut=" + string(nextNut),
	}
}

func writeResponse(w http.ResponseWriter, response *sqrl.ServerMsg) {
	encoded, err := response.Encode()
	if err != nil {
		panic(err)
	}
	// TODO: This is a bit janky but it's what the reference
	// implementation does. Should probably question the use
	// of this content type given it's not in the form key=value.
	w.Header().Set("Content-Type", xFormURLEncoded)
	if _, err := w.Write([]byte(encoded)); err != nil {
		panic(err) // TODO: What to do here?
	}
}
