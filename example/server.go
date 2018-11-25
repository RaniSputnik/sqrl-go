package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	// TODO having to use text/template as html/template mangles
	// the sqrl:// URL and instead turns it into a hash.
	"text/template"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/RaniSputnik/sqrl-go/sqrlhttp"
	"github.com/gorilla/mux"
)

const indexTemplateString = `
<h1>Login With SQRL</h1>
<a href="{{ .LoginURL }}" target="_blank"><img src="data:image/png;base64,{{ .LoginQRCode }}" alt="SQRL Login" /></a>
`

func main() {
	port := 8080

	insecureKey := make([]byte, 16) // TODO: Generate a random key
	sqrlServer := sqrl.Configure(insecureKey)

	d := &delegate{}
	router := mux.NewRouter()
	router.HandleFunc("/", handleIssueChallenge(sqrlServer)).Methods(http.MethodGet)
	router.Handle("/sqrl", sqrlhttp.Authenticate(sqrlServer, d)).Methods(http.MethodPost)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	log.Printf("Server now listening on port: %d", port)
	log.Fatal(server.ListenAndServe())
}

func handleIssueChallenge(server *sqrl.Server) http.HandlerFunc {
	indexTemplate := template.Must(template.New("index").Parse(indexTemplateString))

	return func(w http.ResponseWriter, r *http.Request) {
		loginURL, qrCode := sqrlhttp.GenerateChallenge(server, r, "localhost:8080")
		data := struct {
			LoginURL    string
			LoginQRCode string
		}{
			LoginURL:    loginURL,
			LoginQRCode: base64.StdEncoding.EncodeToString(qrCode),
		}

		indexTemplate.Execute(w, data)
	}
}
