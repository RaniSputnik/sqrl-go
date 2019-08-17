package ssp

import (
	"net/http"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

func (server *Server) PagHandler(store TransactionStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nut := r.URL.Query().Get("nut")

		firstTransaction, err := store.GetFirstTransaction(r.Context(), sqrl.Nut(nut))
		if err != nil {
			// TODO: handle err in some more clever way
			panic(err)
		}

		if firstTransaction == nil || firstTransaction.ClientIP != ClientIP(r) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

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
		_, _ = w.Write([]byte(url))
	})
}
