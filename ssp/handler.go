package ssp

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	sqrl "github.com/RaniSputnik/sqrl-go"
	qrcode "github.com/skip2/go-qrcode"

	"github.com/gorilla/mux"
)

// Handler returns a gorilla mux router including all of the
// SQRL SSP API handlers
func (s *Server) Handler() http.Handler {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/nut.sqrl", s.NutHandler)
	r.HandleFunc("/qr.png", s.QRCodeHandler)
	r.Handle("/cli.sqrl", s.ClientHandler(s.store, s.exchange))
	r.Handle("/pag.sqrl", s.PagHandler(s.store))

	r.Handle("/token", s.TokenHandler(s.exchange)).Methods(http.MethodGet)
	// r.Handle("/users", protect(AddUserHandler(userStore, logger))).Methods(http.MethodPost)
	// r.Handle("/users", protecte(DeleteUserHandler(userStore, logger))).Methods(http.MethodDelete)

	return r
}

// NutHandler handler for the nut endpoint
// Reference: https://www.grc.com/sqrl/sspapi.htm
// TODO does not yet handle params 0-9, sin or ask
func (s *Server) NutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	nut := s.Nut(clientID(r))
	s.logger.Printf("Generated nut: %s", nut)

	formValues := make(url.Values)
	formValues.Add("nut", string(nut))
	formValues.Add("can", sqrl.Base64.EncodeToString([]byte(r.Header.Get("Referer"))))

	if _, err := w.Write([]byte(formValues.Encode())); err != nil {
		s.logger.Printf("Nut write unsuccessful: %v", err)
	}
}

// QRCodeHandler handles creating the QR code version
// of the SQRL URL
func (s *Server) QRCodeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nut := query.Get("nut")
	size := atoiWithDefault(query.Get("size"), 256)

	if nut == "" {
		s.logger.Printf("QR code requested with empty 'nut' parameter")
		w.WriteHeader(http.StatusBadRequest) // TODO: default pending image
		_, _ = w.Write([]byte("Missing nut param"))
		return
	}
	if r.Host == "" {
		s.logger.Printf("QR code requested with no 'Host' header set")
		w.WriteHeader(http.StatusBadRequest) // TODO: default error image
		_, _ = w.Write([]byte("Missing Host header"))
		return
	}

	params := make(url.Values)
	params.Add("nut", nut)
	loginURL := url.URL{
		Scheme:   "sqrl",
		Host:     requestDomain(r),
		RawQuery: params.Encode(),
	}

	bytes, err := qrcode.Encode(loginURL.String(), qrcode.Medium, size)
	if err != nil {
		s.logger.Printf("Failed to encode login URL '%s': %v", loginURL, err)
		w.WriteHeader(http.StatusNotFound) // TODO: default error image
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		s.logger.Printf("QR code write unsuccessful: %v", err)
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
	if forwardedAddress := r.Header.Get("X-Forwarded-For"); forwardedAddress != "" {
		return forwardedAddress
	}
	return r.RemoteAddr
}

func requestDomain(r *http.Request) string {
	return r.Host
}

func getTokenRedirectURL(server *Server, token Token) string {
	return fmt.Sprintf("%s?token=%s", server.redirectURL, token)
}
