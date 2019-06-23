package ssp

import (
	"time"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

type Server struct {
	key            []byte
	redirectURL    string
	clientEndpoint string

	nutter *sqrl.Server
}

func Configure(key []byte, redirectURL string) *Server {
	sqrlServer := sqrl.Configure(key)

	return &Server{
		key:            key,
		redirectURL:    redirectURL,
		clientEndpoint: "/cli.sqrl",

		// TODO: This should actually be named something
		// more like "nut generator"
		nutter: sqrlServer,
	}
}

func (s *Server) Nut(clientIdentifier string) sqrl.Nut {
	return s.nutter.Nut(clientIdentifier)
}

// WithNutExpiry sets the window of time within which
// a nut is considered to be valid.
func (s *Server) WithNutExpiry(d time.Duration) *Server {
	s.nutter.WithNutExpiry(d)
	return s
}

// WithClientEndpoint sets the endpoint that the client can
// use to post SQRL transactions to. This endpoint should
// be the path relative to the SQRL domain eg. /sqrl/cli.sqrl
//
// Defaults to /cli.sqrl if not set.
func (s *Server) WithClientEndpoint(url string) *Server {
	s.clientEndpoint = url
	return s
}
