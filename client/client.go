package client

import (
	"errors"
	"net/url"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

var (
	ErrUriInvalid = errors.New("uri invalid")
)

func Login(uri string) error {
	parsed, err := url.Parse(uri)
	if err != nil {
		return ErrUriInvalid
	}

	if parsed.Scheme != sqrl.Scheme {
		return ErrUriInvalid
	}

	return errors.New("TODO")
}
