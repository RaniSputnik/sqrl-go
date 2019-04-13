package ssp

import (
	"fmt"
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

		w.WriteHeader(http.StatusFound)
		http.Redirect(w, r, url, http.StatusFound)
	})
}
