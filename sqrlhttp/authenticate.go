package sqrlhttp

import (
	"log"
	"net/http"

	"github.com/RaniSputnik/sqrl-go"
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
		log.Printf("Got SQRL request: %v", r)

		response := genNextResponse(server, r)
		defer writeResponse(w, response)

		if r.Header.Get("Content-Type") != xFormURLEncoded {
			clientFailure(response)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("Failed to parse form: %s", err)
			clientFailure(response)
			return
		}

		clientRaw := r.Form.Get("client")
		serverRaw := r.Form.Get("server")

		client, errc := sqrl.ParseClient(clientRaw)
		// Ensure that it is the URL / params we sent
		//_, errs := sqrl.ParseServer(r.Form.Get("server"))
		if errc != nil /*|| errs != nil */ {
			clientFailure(response)
			return
		}

		// TODO: Verify server parameter in some way

		ids := sqrl.Signature(r.Form.Get("ids"))
		signedPayload := clientRaw + serverRaw
		if !ids.Verify(client.Idk, signedPayload) {
			clientFailure(response)
			return
		}

		// TODO: Test for IP Match

		isKnown, err := delegate.Known(r.Context(), client.Idk)
		if err != nil {
			log.Printf("Failed to process response: %v", err)
			serverError(response)
			return
		} else if isKnown {
			response.Tif |= sqrl.TIFCurrentIDMatch
		}
	})
}

func genNextResponse(server *sqrl.Server, r *http.Request) *sqrl.ServerMsg {
	nextNut := server.Nut(clientID(r))
	return &sqrl.ServerMsg{
		Ver: v1Only,
		Nut: nextNut,
		Qry: r.URL.Path + "?nut=" + string(nextNut),
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
	w.Write([]byte(encoded))
}
