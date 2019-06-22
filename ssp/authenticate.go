package ssp

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

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

func Authenticate(server *sqrl.Server, store Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got SQRL request: %v\n", r)

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
		if errc != nil {
			log.Printf("Client param (%s) invalid: %v", clientRaw, errc)
			clientFailure(response)
			return
		}
		nut, serverOK := verifyServer(serverRaw)
		if !serverOK {
			log.Printf("Server param (%s) invalid", serverRaw)
			clientFailure(response)
			return
		}
		// TODO: Verify nut hasn't expired

		signedPayload := clientRaw + serverRaw
		if !ids.Verify(client.Idk, signedPayload) {
			clientFailure(response)
			return
		}
		// TODO: Verify previous identity signatures

		thisTransaction := &Transaction{
			Id:   nut,
			Next: response.Nut,
		}
		if err := store.SaveTransaction(r.Context(), thisTransaction); err != nil {
			log.Printf("Failed to save transaction: %v\n", err)
			serverError(response)
			return
		}
		firstTransaction, err := store.GetFirstTransaction(r.Context(), nut)
		if err != nil {
			log.Printf("Failed to retrieve first transaction: %v\n", err)
			serverError(response)
			return
		}
		if firstTransaction == nil {
			firstTransaction = thisTransaction
		}

		// TODO: Fetch user details
		// TODO: Fetch user from previous identity

		// TODO: Test for IP Match

		// TODO: Pass previous identities to "known"
		isKnown, err := store.GetIsKnown(r.Context(), client.Idk)
		if err != nil {
			log.Printf("Failed to determine if identity is known: %v\n", err)
			serverError(response)
			return
		} else if isKnown {
			response.Tif |= sqrl.TIFCurrentIDMatch
		}

		switch client.Cmd {
		case sqrl.CmdIdent:
			token := "todo-token"
			err := store.SaveIdentSuccess(r.Context(), firstTransaction.Id, token)
			if err != nil {
				log.Fatalf("Failed to check authenticated: %v\n", err)
				serverError(response)
			}

			if client.HasOpt(sqrl.OptCPS) {
				response.URL = fmt.Sprintf("%s?%s", server.RedirectURL(), token)
			}
		case sqrl.CmdQuery:
			// TODO: Anything need to be done here?

		default:
			// In all other cases, not supported
			response.Tif |= sqrl.TIFFunctionNotSupported
		}
	})
}

func genNextResponse(server *sqrl.Server, r *http.Request) *sqrl.ServerMsg {
	nextNut := server.Nut(clientID(r))
	return &sqrl.ServerMsg{
		Ver: v1Only,
		Nut: nextNut,
		Qry: server.ClientEndpoint() + "?nut=" + string(nextNut),
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

func verifyServer(serverRaw string) (sqrl.Nut, bool) {
	// TODO: Here we accept EITHER a URL or ServerMsg
	// However we know that ONLY the first request
	// from the client should be a URL.
	// Is there a way for us to ensure that here?

	bytes, err := sqrl.Base64.DecodeString(serverRaw)
	if err != nil {
		return "", false
	}

	server := string(bytes)
	if strings.HasPrefix(server, "sqrl") {
		serverURL, err := url.Parse(server)
		if err != nil {
			return "", false
		}

		// TODO: Assert URL matches server configuration
		// eg. domain, "server friendly name", etc.

		nut := serverURL.Query().Get("nut")
		if nut == "" {
			return "", false
		}
		return sqrl.Nut(nut), true

	} else {
		msg, err := sqrl.ParseServer(serverRaw)
		if err != nil || msg == nil || msg.Nut == "" {
			return "", false
		}
		return msg.Nut, true
	}
}
