package ssp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	sqrl "github.com/RaniSputnik/sqrl-go"
	qrcode "github.com/skip2/go-qrcode"

	"github.com/gorilla/mux"
)

func Handler(s *sqrl.Server, authFunc ServerToServerAuthValidationFunc) http.Handler {
	// TODO: Make this configurable
	logger := log.New(os.Stdout, "", 0)

	store := NewMemoryStore()

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/nut.json", nutHandler(s, logger))
	r.HandleFunc("/qr.png", qrHandler(s, logger))
	r.Handle("/cli.sqrl", Authenticate(s, store))
	r.Handle("/pag.sqrl", PagHandler(s, store))

	userStore := &todoUserStore{}
	protect := ServerToServerAuthMiddleware(authFunc, logger)
	r.Handle("/token", protect(TokenHandler(s, userStore, logger))).Methods(http.MethodGet)
	// r.Handle("/users", protect(AddUserHandler(userStore, logger))).Methods(http.MethodPost)
	// r.Handle("/users", protecte(DeleteUserHandler(userStore, logger))).Methods(http.MethodDelete)

	return r
}

func nutHandler(server *sqrl.Server, logger *log.Logger) http.HandlerFunc {
	type nutResponse struct {
		Nut string `json:"nut"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		res := nutResponse{
			Nut: server.Nut(clientID(r)).String(),
		}
		logger.Printf("Generated nut: %s", res.Nut)

		if err := json.NewEncoder(w).Encode(res); err != nil {
			logger.Printf("Nut write unsuccessful: %v", err)
		}
	}
}

func qrHandler(server *sqrl.Server, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		nut := query.Get("nut")
		size := atoiWithDefault(query.Get("size"), 256)

		if nut == "" {
			logger.Printf("QR code requested with empty 'nut' parameter")
			w.WriteHeader(http.StatusNotFound) // TODO: default pending image
			return
		}
		if r.Host == "" {
			logger.Printf("QR code requested with no 'Host' header set")
			w.WriteHeader(http.StatusNotFound) // TODO: default error image
			return
		}

		loginURL := fmt.Sprintf("sqrl://%s/sqrl?nut=%s", requestDomain(r), nut)
		bytes, err := qrcode.Encode(loginURL, qrcode.Medium, size)
		if err != nil {
			logger.Printf("Failed to encode login URL '%s': %v", loginURL, err)
			w.WriteHeader(http.StatusNotFound) // TODO: default error image
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(bytes); err != nil {
			logger.Printf("QR code write unsuccessful: %v", err)
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

func getTokenRedirectURL(server *sqrl.Server, token string) string {
	return fmt.Sprintf("%s?token=%s", server.RedirectURL(), token)
}
