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
	"github.com/gorilla/mux"
	qrcode "github.com/skip2/go-qrcode"
)

const indexTemplateString = `
<h1>Login With SQRL</h1>
<a href="{{ .LoginURL }}" target="_blank"><img src="data:image/png;base64,{{ .LoginQRCode }}" alt="SQRL Login" /></a>
`

func main() {
	port := 8080

	indexTemplate := template.Must(template.New("index").Parse(indexTemplateString))

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loginURL := createLoginURL(r, "localhost:8080")
		qrCode, _ := createQRCode(loginURL) // TODO handle error case

		data := struct {
			LoginURL    string
			LoginQRCode string
		}{
			LoginURL:    loginURL,
			LoginQRCode: qrCode,
		}

		indexTemplate.Execute(w, data)
	})

	// TODO handle callback URL

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	log.Printf("Server now listening on port: %d", port)
	log.Fatal(server.ListenAndServe())
}

func createLoginURL(r *http.Request, domain string) string {
	nonce := sqrl.Nut(r)
	return fmt.Sprintf("sqrl://%s/sqrl?%s", domain, nonce)
}

func createQRCode(URL string) (string, error) {
	png, err := qrcode.Encode(URL, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(png), nil
}
