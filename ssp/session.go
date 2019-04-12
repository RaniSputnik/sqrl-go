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

		id, state, err := delegate.GetSession(r.Context(), sqrl.Nut(nut))
		if err != nil {
			// TODO: handle err in some more clever way
			panic(err)
		}

		switch state {
		case SessionUnknown:
			fallthrough
		case SessionAuthenticating:
			w.WriteHeader(http.StatusNotFound)
		case SessionAuthenticated:
			// TODO: We should probably save the token in the ident handler
			// so that it can be verified by the resource server later
			token := "todo-token"
			url := fmt.Sprintf("%s?%s", server.RedirectURL(), token)

			if err := delegate.Redirected(r.Context(), id, token); err != nil {
				panic(err) // TODO: Handle error in some sensible way
			}

			w.WriteHeader(http.StatusFound)
			http.Redirect(w, r, url, http.StatusFound)
		}
	})
}
