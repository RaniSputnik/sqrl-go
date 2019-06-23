package ssp

import (
	"time"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type donothingLogger struct{}

func (_ donothingLogger) Printf(format string, v ...interface{}) {}

type Server struct {
	key []byte

	logger         Logger
	redirectURL    string
	clientEndpoint string

	validator ServerToServerAuthValidationFunc

	nutter *sqrl.Server
}

func Configure(key []byte, redirectURL string) *Server {
	sqrlServer := sqrl.Configure(key)

	return &Server{
		key: key,

		logger:         donothingLogger{},
		redirectURL:    redirectURL,
		clientEndpoint: "/cli.sqrl",

		// TODO: Is there a more sensible default we could use here?
		validator: noProtection,

		// TODO: This should actually be named something
		// more like "nut generator"
		nutter: sqrlServer,
	}
}

func (s *Server) WithLogger(l Logger) *Server {
	if l == nil {
		l = donothingLogger{}
	}
	s.logger = l
	return s
}

func (s *Server) WithAuthentication(validator ServerToServerAuthValidationFunc) *Server {
	s.validator = validator
	return s
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

func (s *Server) Nut(clientIdentifier string) sqrl.Nut {
	return s.nutter.Nut(clientIdentifier)
}
