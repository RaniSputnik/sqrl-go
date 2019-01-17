package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/RaniSputnik/sqrl-go/ssp"
)

// TODO: Do not use this
var todoKey = make([]byte, 16)

func main() {
	dir := "static"
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	// TODO: Don't strip the trailing slash here or else gorilla Mux will become confused
	// and attempt to clean+rediect. Is this something that we should handle in library code?
	http.Handle("/sqrl/", http.StripPrefix("/sqrl", ssp.Handler(todoKey)))
	http.Handle("/", indexHandler())

	port := ":8080"
	log.Printf("Serving files from './%s' on port %s", dir, port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
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
