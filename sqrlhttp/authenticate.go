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

func Authenticate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextNut := sqrl.Nut(r)
		response := &sqrl.ServerMsg{
			Ver: v1Only,
			Nut: nextNut,
			Qry: r.URL.Path + "?nut=" + nextNut,
		}

		// Write the response
		defer func() {
			encoded, err := response.Encode()
			if err != nil {
				panic(err)
			}
			// TODO: This is a bit janky but it's what the reference
			// implementation does. Should probably question the use
			// of this content type given it's not in the form key=value.
			w.Header().Set("Content-Type", xFormURLEncoded)
			w.Write([]byte(encoded))
		}()

		if r.Header.Get("Content-Type") != xFormURLEncoded {
			clientFailure(response)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("Failed to parse form: %s", err)
			clientFailure(response)
			return
		}

		_, errc := sqrl.ParseClient(r.Form.Get("client"))
		_, errs := sqrl.ParseServer(r.Form.Get("server"))
		if errc != nil || errs != nil {
			clientFailure(response)
			return
		}

		// TODO: how do we parse serverMsg (and do we need to?)
		//serverMsg, errs := sqrl.ParseServer(r.Form.Get("server"))

		// TODO: Verify signatures

		// TODO: Test for IP Match
	})
}
