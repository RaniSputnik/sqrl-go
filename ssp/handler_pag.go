package ssp

import (
	"log"
	"net/http"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

func (server *Server) PagHandler(store TransactionStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nut := r.URL.Query().Get("nut")

		// TODO: Do not call the store with an invalid nut

		// TODO: Check for an IP match (IP should ALWAYS match here)

		token, err := store.GetIdentSuccess(r.Context(), sqrl.Nut(nut))
		if err != nil {
			// TODO: handle err in some more clever way
			panic(err)
		}

		if token == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		url := getTokenRedirectURL(server, token)
		if _, err := w.Write([]byte(url)); err != nil {
			log.Printf("Failed to write pag.sqrl response: %v", err)
		}
	})
}
