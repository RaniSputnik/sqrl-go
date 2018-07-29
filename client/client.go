package client

import (
	"crypto"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"golang.org/x/crypto/ed25519"
)

var (
	ErrUriInvalid = errors.New("uri invalid")
)

var HttpClient = &http.Client{
	Timeout: time.Second * 5,
}

func Login(uri string) error {
	var rand io.Reader
	pub, priv, err := ed25519.GenerateKey(rand)
	if err != nil {
		return err
	}

	endpoint, err := getEndpoint(uri)
	if err != nil {
		return err
	}

	idk := sqrl.Identity(sqrl.Base64.EncodeToString(pub))

	clientParameters := QueryCmd(idk)
	serverParameters := sqrl.Base64.EncodeToString([]byte(uri))

	signMe := clientParameters + serverParameters

	ids, err := sign(signMe, priv)
	if err != nil {
		return err
	}

	form := []string{
		"client=" + clientParameters,
		"server=" + serverParameters,
		"ids=" + ids,
	}
	_, err = do(endpoint, strings.Join(form, "&"))
	if err != nil {
		return err
	}

	return nil
}

func do(uri string, form string) (*sqrl.ServerMsg, error) {
	res, err := HttpClient.Post(uri, "application/x-www-form-urlencoded", strings.NewReader(form))
	if err != nil {
		return nil, err
	}

	gotBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return sqrl.ParseServer(string(gotBody))
}

// getEndpoint transforms a sqrl:// URL to a https:// URL
func getEndpoint(uri string) (endpoint string, err error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return "", ErrUriInvalid
	}
	if parsed.Scheme != sqrl.Scheme {
		return "", ErrUriInvalid
	}

	parsed.Scheme = "https"
	return parsed.String(), nil
}

// sign accepts a payload to sign with the given private key
//
// The payload should be the value of the 'server' parameter
// appended to the value of the 'client' parameter.
func sign(payload string, privateKey ed25519.PrivateKey) (string, error) {
	fmt.Printf("Signing: '%s'\n", payload)
	sig, err := privateKey.Sign(nil, []byte(payload), crypto.Hash(0))
	if err != nil {
		return "", err
	}
	result := sqrl.Base64.EncodeToString(sig)
	fmt.Printf("Signature: '%s'\n", result)
	return result, nil
}
