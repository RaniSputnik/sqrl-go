package ssp

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	sqrl "github.com/RaniSputnik/sqrl-go"

	"github.com/gorilla/mux"
)

func Handler(key []byte) http.Handler {
	// TODO: Make this configurable
	logger := log.New(os.Stdout, "", 0)
	s := sqrl.Configure(key)

	// TODO: Why does this redirect when used with StripPrefix?
	r := mux.NewRouter().StrictSlash(false)

	r.HandleFunc("/nut.json", nutHandler(s, logger))
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
			logger.Printf("Write unsuccessful: %v", err)
		}
	}
}
