package main

import (
	"fmt"
	"log"
	"net/http"

	"html/template"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/RaniSputnik/sqrl-go/ssp"
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

var logins = &challengeStore{
	vals: make(map[string]LoginState),
}

func main() {
	port := 8080

	insecureKey := make([]byte, 16) // TODO: Generate a random key
	sqrlServer := sqrl.Configure(insecureKey)

	d := &delegate{}
	router := mux.NewRouter()
	router.HandleFunc("/", handleIndex()).Methods(http.MethodGet)
	router.HandleFunc("/login", handleIssueChallenge(sqrlServer)).Methods(http.MethodGet)
	router.HandleFunc("/logout", handleLogout()).Methods(http.MethodGet)
	router.Handle("/sqrl", ssp.Authenticate(sqrlServer, d)).Methods(http.MethodPost)
	router.HandleFunc("/sync.txt", handleSync()).Methods(http.MethodGet)

	// TODO: Remove this once we have SQRL login working
	// For now this is just for testing sync
	router.HandleFunc("/login/bypass", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "sess")
		challenge := session.Values["challenge"].(string)
		logins.Set(challenge, LoginStateAuthenticated)
	})

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
	notAuthenticated := []byte("not authenticated")

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "sess")
		challenge, ok := session.Values["challenge"].(string)
		if !ok {
			w.Write(notAuthenticated)
			return
		}

		fmt.Printf("Sync: Retrieving challenge %s\n", challenge)
		loginState := logins.Get(challenge)

		fmt.Printf("Sync: Current login state = %s\n", loginState)
		if loginState == LoginStateAuthenticated {
			// TODO: Check IP still matches the original challenge
			session.Values["uid"] = "TODO: Set the user ID"
			session.Save(r, w)
			w.Write([]byte("authenticated"))
			fmt.Println("Sync: Authentication successfull!")
			return
		}

		w.Write(notAuthenticated)
	}
}

func handleIssueChallenge(server *sqrl.Server) http.HandlerFunc {
	loginTemplate := template.Must(template.New("login").Parse(loginTemplateString))

	return func(w http.ResponseWriter, r *http.Request) {
		challenge, loginFragment := ssp.GenerateChallenge(server, r)
		must(logins.Set(challenge, LoginStateIssued))
		fmt.Printf("Issue: New challenge = %s", challenge)

		session, _ := store.Get(r, "sess")
		session.Values["challenge"] = challenge
		session.Save(r, w)

		data := struct{ SQRL template.HTML }{SQRL: loginFragment}
		loginTemplate.Execute(w, data)
	}
}

func handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "sess")
		delete(session.Values, "challenge")
		delete(session.Values, "uid")
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
