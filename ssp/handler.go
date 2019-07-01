package ssp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	qrcode "github.com/skip2/go-qrcode"

	"github.com/gorilla/mux"
)

func (s *Server) Handler() http.Handler {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/nut.json", s.NutHandler())
	r.HandleFunc("/qr.png", s.QRCodeHandler())
	r.Handle("/cli.sqrl", s.ClientHandler(s.store, s.exchange))
	r.Handle("/pag.sqrl", s.PagHandler(s.store))

	r.Handle("/token", s.TokenHandler(s.exchange)).Methods(http.MethodGet)
	// r.Handle("/users", protect(AddUserHandler(userStore, logger))).Methods(http.MethodPost)
	// r.Handle("/users", protecte(DeleteUserHandler(userStore, logger))).Methods(http.MethodDelete)

	return r
}

func (server *Server) NutHandler() http.HandlerFunc {
	type nutResponse struct {
		Nut string `json:"nut"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		res := nutResponse{
			Nut: server.Nut(clientID(r)).String(),
		}
		server.logger.Printf("Generated nut: %s", res.Nut)

		if err := json.NewEncoder(w).Encode(res); err != nil {
			server.logger.Printf("Nut write unsuccessful: %v", err)
		}
	}
}

func (server *Server) QRCodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		nut := query.Get("nut")
		size := atoiWithDefault(query.Get("size"), 256)

		if nut == "" {
			server.logger.Printf("QR code requested with empty 'nut' parameter")
			w.WriteHeader(http.StatusNotFound) // TODO: default pending image
			return
		}
		if r.Host == "" {
			server.logger.Printf("QR code requested with no 'Host' header set")
			w.WriteHeader(http.StatusNotFound) // TODO: default error image
			return
		}

		loginURL := fmt.Sprintf("sqrl://%s/sqrl?nut=%s", requestDomain(r), nut)
		bytes, err := qrcode.Encode(loginURL, qrcode.Medium, size)
		if err != nil {
			server.logger.Printf("Failed to encode login URL '%s': %v", loginURL, err)
			w.WriteHeader(http.StatusNotFound) // TODO: default error image
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(bytes); err != nil {
			server.logger.Printf("QR code write unsuccessful: %v", err)
		}
	}
}

// helpers

func atoiWithDefault(val string, def int) int {
	if res, err := strconv.Atoi(val); err != nil {
		return def
	} else {
		return res
	}
}

func clientID(r *http.Request) string {
	// TODO: X-Forwarded-For
	// TODO: Include user agent if available
	return r.RemoteAddr
}

func requestDomain(r *http.Request) string {
	// TODO: Do we need to do anything special here for proxies?
	return r.Host
}

func getTokenRedirectURL(server *Server, token Token) string {
	return fmt.Sprintf("%s?token=%s", server.redirectURL, token)
}
