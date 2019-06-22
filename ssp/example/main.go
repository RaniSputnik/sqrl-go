package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/RaniSputnik/sqrl-go/ssp"
)

// TODO: Do not use this
var todoKey = make([]byte, 16)

const clientSecret = "something-very-secret"

func main() {
	// TODO: This builder is a bit gross
	// Maybe we can move to using option functions
	// like Gorilla Handlers?
	// http://www.gorillatoolkit.org/pkg/handlers#CORSOption
	config := sqrl.Configure(todoKey).
		WithNutExpiry(time.Minute * 5).
		WithRedirectURL("http://localhost:8080/callback").
		// TODO: bit lame that this cli.sqrl is both hardcoded
		// in ssp and configured here. Should we only provide
		// the /sqrl part here? Or should cli.sqrl be moved out
		// of ssp.Handler?
		WithCLientEndpoint("/sqrl/cli.sqrl")

	serverToServerProtection := func(r *http.Request) error {
		if r.Header.Get("X-Client-Secret") != clientSecret {
			return errors.New("Invalid X-Client-Secret header")
		}
		return nil
	}

	dir := "static"
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	// TODO: Don't strip the trailing slash here or else gorilla Mux will become confused
	// and attempt to clean+rediect. Is this something that we should handle in library code?
	http.Handle("/sqrl/", http.StripPrefix("/sqrl", ssp.Handler(config, serverToServerProtection)))
	http.Handle("/callback", authCallbackHandler())
	http.Handle("/", indexHandler())

	port := ":8080"
	log.Printf("Serving files from './%s' on port %s", dir, port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func authCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			log.Printf("Callback called without token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := w.Write([]byte("TODO")); err != nil {
			log.Printf("Failed to write callback response: %v", err)
		}
	}
}

func indexHandler() http.HandlerFunc {
	type templateData struct {
		Session string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: This parsing should not be done in the http.Handler
		// I have left it here for easy development.
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, templateData{
			Session: "TODO",
		}); err != nil {
			log.Printf("Failed to render template: %v", err)
		}
	}
}
