package main

import (
	"fmt"
	"log"
	"net/http"

	"html/template"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/RaniSputnik/sqrl-go/sqrlhttp"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

const indexTemplateString = `
<h1>Welcome {{ .UserID }}</h1>
`
const loginTemplateString = `
<h1>Login With SQRL</h1>
{{ .SQRL }}
`

func main() {
	port := 8080

	insecureKey := make([]byte, 16) // TODO: Generate a random key
	sqrlServer := sqrl.Configure(insecureKey)

	d := &delegate{}
	router := mux.NewRouter()
	router.HandleFunc("/", handleIndex()).Methods(http.MethodGet)
	router.HandleFunc("/login", handleIssueChallenge(sqrlServer)).Methods(http.MethodGet)
	router.Handle("/sqrl", sqrlhttp.Authenticate(sqrlServer, d)).Methods(http.MethodPost)
	router.HandleFunc("/sync.txt", handleSync()).Methods(http.MethodGet)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	log.Printf("Server now listening on port: %d", port)
	log.Fatal(server.ListenAndServe())
}

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

func handleIndex() http.HandlerFunc {
	indexTemplate := template.Must(template.New("index").Parse(indexTemplateString))

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "sess")
		uid, authenticated := session.Values["uid"].(string)
		if !authenticated {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		data := struct{ UserID string }{UserID: uid}
		indexTemplate.Execute(w, data)
	}
}

// TODO: How might we move this into the SQRL library?
func handleSync() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: check the users outstanding SQRL request
		// If the users SQRL request has been accepted
		// Then set an 'authenticated' cookie and return success
		// The user will now be redirected to login

		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
	}
}

func handleIssueChallenge(server *sqrl.Server) http.HandlerFunc {
	loginTemplate := template.Must(template.New("login").Parse(loginTemplateString))

	return func(w http.ResponseWriter, r *http.Request) {
		loginFragment := sqrlhttp.GenerateChallenge(server, r, "localhost:8080")
		data := struct{ SQRL template.HTML }{SQRL: loginFragment}
		loginTemplate.Execute(w, data)
	}
}
