package ssp

import (
	"fmt"
	"log"
	"net/http"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

type SessionState string

const (
	SessionUnknown        = SessionState("unknown")
	SessionAuthenticating = SessionState("authenticating")
	SessionAuthenticated  = SessionState("authenticated")
)

func SessionHandler(server *sqrl.Server, delegate Delegate) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nut := r.URL.Query().Get("nut")

		// TODO: Do not call the delgate with an invalid nut

		// TODO: Check for an IP match (IP should ALWAYS match here)

		id, token, err := delegate.GetSession(r.Context(), sqrl.Nut(nut))
		if err != nil {
			// TODO: handle err in some more clever way
			panic(err)
		}

		// TODO: If we have an ID

		if token == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		url := fmt.Sprintf("%s?%s", server.RedirectURL(), token)
		if err := delegate.Redirected(r.Context(), id); err != nil {
			panic(err) // TODO: Handle error in some sensible way
		}

		if _, err := w.Write([]byte(url)); err != nil {
			log.Printf("Failed to write pag.sqrl response: %v", err)
		}
	})
}
