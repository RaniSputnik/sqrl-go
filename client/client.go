package client

import (
	"crypto"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/cretz/bine/torutil/ed25519"
)

var (
	ErrUriInvalid = errors.New("uri invalid")
)

var HttpClient = &http.Client{
	Timeout: time.Second * 5,
}

func Login(uri string) error {
	var rand io.Reader
	keypair, err := ed25519.GenerateKey(rand)
	if err != nil {
		return err
	}
	pub := keypair.PublicKey()
	priv := keypair.PrivateKey()

	parsed, err := url.Parse(uri)
	if err != nil {
		return ErrUriInvalid
	}

	if parsed.Scheme != sqrl.Scheme {
		return ErrUriInvalid
	}

	parsed.Scheme = "https"

	idk := b64(string(pub))
	clientParameters := QueryCmd(idk)

	vals := strings.Split(uri, "?")
	vals[1] = url.QueryEscape(vals[1])
	encodedURI := strings.Join(vals, "?")
	fmt.Printf("Encoded server URI:\nFrom '%s'\nTo   '%s'\n", uri, encodedURI)
	serverParameters := b64(encodedURI)

	ids, err := sign(clientParameters, serverParameters, priv)
	if err != nil {
		return err
	}
	if !verify(idk, ids, clientParameters+serverParameters) {
		fmt.Printf("Uh oh, verification failed... wonder why...")
		return errors.New("Sanity check failed")
	}

	form := []string{
		"client=" + clientParameters,
		"server=" + serverParameters,
		"ids=" + ids,
	}
	body := strings.NewReader(strings.Join(form, "&"))

	_, err = HttpClient.Post(parsed.String(), "application/x-www-form-urlencoded", body)
	if err != nil {
		return err
	}

	return nil
}

func sign(client, server string, privateKey ed25519.PrivateKey) (string, error) {
	valueToSign := client + server
	fmt.Printf("Signing: '%s'\n", valueToSign)

	sig, err := privateKey.Sign(nil, []byte(valueToSign), crypto.Hash(0))
	if err != nil {
		return "", err
	}
	result := b64(string(sig))
	fmt.Printf("Signature: '%s'\n", result)
	return result, nil
}

func verify(idk, ids, params string) bool {
	decoder := base64.URLEncoding.WithPadding(base64.NoPadding)
	idkDecoded, err1 := decoder.DecodeString(idk)
	idsDecoded, err2 := decoder.DecodeString(ids)
	if err1 != nil || err2 != nil {
		fmt.Printf("Failed to verify:\ninput='%s', err='%v'\ninput='%s', err2='%v'\n",
			idk, err1, ids, err2)
		return false
	}
	pub := ed25519.PublicKey([]byte(idkDecoded))
	return ed25519.Verify(pub, []byte(params), idsDecoded)
}
