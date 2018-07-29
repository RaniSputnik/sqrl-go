package sqlhttp

import (
	"log"
	"net/http"

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

		server, errs := sqrl.Base64.DecodeString(r.Form.Get("server"))
		client, errc := sqrl.Base64.DecodeString(r.Form.Get("client"))
		if errc != nil || errs != nil || len(client) == 0 || len(server) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//clientVals := strings.Split(client, "\n")
	})
}
