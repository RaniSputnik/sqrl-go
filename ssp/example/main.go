package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RaniSputnik/sqrl-go/ssp"
)

// TODO: Do not use this
var todoKey = make([]byte, 16)

const clientSecret = "something-very-secret"

const (
	CookieSQRLUser = "sqrl_user"
)

func main() {
	serverToServerProtection := func(r *http.Request) error {
		if r.Header.Get("X-Client-Secret") != clientSecret {
			return errors.New("Invalid X-Client-Secret header")
		}
		return nil
	}

	// TODO: This builder is a bit gross
	// Maybe we can move to using option functions
	// like Gorilla Handlers?
	// http://www.gorillatoolkit.org/pkg/handlers#CORSOption
	sspServer := ssp.Configure(todoKey, "http://localhost:8080/callback").
		WithAuthentication(serverToServerProtection).
		WithLogger(log.New(os.Stdout, "SSP: ", 0)).
		WithNutExpiry(time.Minute * 5).
		// TODO: bit lame that this cli.sqrl is both hardcoded
		// in ssp and configured here. Should we only provide
		// the /sqrl part here? Or should cli.sqrl be moved out
		// of ssp.Handler?
		WithClientEndpoint("/sqrl/cli.sqrl")

	dir := "static"
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	// TODO: Don't strip the trailing slash here or else gorilla Mux will become confused
	// and attempt to clean+rediect. Is this something that we should handle in library code?
	http.Handle("/sqrl/", http.StripPrefix("/sqrl", sspServer.Handler()))
	http.Handle("/callback", authCallbackHandler("http://localhost:8080/sqrl/token"))
	http.Handle("/logout", logoutHandler())
	http.Handle("/", indexHandler())

	port := ":8080"
	log.Printf("Serving files from './%s' on port %s", dir, port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func authCallbackHandler(sspTokenURL string) http.HandlerFunc {
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			log.Printf("Callback called without token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userId, err := validateToken(client, sspTokenURL, token)
		if err != nil {
			log.Printf("Failed to validate token: %+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if userId == "" {
			log.Printf("Invalid token: %s", token)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Note: this is not a secure way of setting cookies
		// for authentication state, you should use something
		// like http://www.gorillatoolkit.org/pkg/securecookie
		// in a real project.
		http.SetCookie(w, &http.Cookie{
			Name:  CookieSQRLUser,
			Value: userId,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(CookieSQRLUser); err == nil {
			cookie.Expires = time.Now()
			http.SetCookie(w, cookie)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func indexHandler() http.HandlerFunc {
	type templateData struct {
		Authenticated bool
		UserID        string
		Session       string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: This parsing should not be done in the http.Handler
		// I have left it here for easy development.
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var userID string
		if cookie, err := r.Cookie(CookieSQRLUser); err == nil {
			userID = cookie.Value
		}

		if err := tmpl.Execute(w, templateData{
			Authenticated: userID != "",
			UserID:        userID,
			Session:       "TODO",
		}); err != nil {
			log.Printf("Failed to render template: %v", err)
		}
	}
}

func validateToken(client *http.Client, sspTokenURL string, token string) (userId string, err error) {
	type tokenResponse struct {
		User string `json:"user"`
	}

	url := fmt.Sprintf("%s?token=%s", sspTokenURL, token)
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	r.Header.Set("X-Client-Secret", clientSecret)
	res, err := client.Do(r)
	if err != nil {
		return "", err
	}
	if res.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("remote server returned error when validating token: %s", res.Status)
	}

	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)

	var got tokenResponse
	err = json.Unmarshal(bytes, &got)
	if err != nil {
		log.Printf("Failed to decode body: %s", bytes)
	}
	return got.User, err
}
