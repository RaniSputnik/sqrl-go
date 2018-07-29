package http

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/RaniSputnik/sqrl-go"
)

const xFormURLEncoded = "application/x-www-form-urlencoded"

type request struct {
}

type clientParams struct {
	Ver string
	Cmd sqrl.Cmd
}

func Authenticate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.UserAgent() != V1 || r.Header.Get("Content-Type") != xFormURLEncoded {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("Failed to parse form: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		server, errs := b64decode(r.Form.Get("server"))
		client, errc := b64decode(r.Form.Get("client"))
		if errc != nil || errs != nil || client == "" || server == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		clientVals := strings.Split(client, "\n")
	})
}

func b64decode(in string) (string, error) {
	got, err := base64.StdEncoding.DecodeString(in)
	return string(got), err
}
