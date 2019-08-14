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

	store          Store
	exchange       TokenExchange
	logger         Logger
	validator      ServerToServerAuthValidationFunc
	redirectURL    string
	clientEndpoint string

	nutter sqrl.Nutter
}

func Configure(key []byte, redirectURL string) *Server {
	store := NewMemoryStore()
	exchange := DefaultExchange(key, time.Minute)
	nutter := sqrl.NewNutter(key)

	return &Server{
		key: key,

		store:    store,
		exchange: exchange,
		logger:   donothingLogger{},
		// TODO: Is there a more sensible default we could use here?
		validator:      noProtection,
		redirectURL:    redirectURL,
		clientEndpoint: "/cli.sqrl",

		nutter: nutter,
	}
}

func (s *Server) WithStore(store Store) *Server {
	s.store = store
	return s
}

func (s *Server) WithTokenExchange(exchange TokenExchange) *Server {
	s.exchange = exchange
	return s
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
	// TODO: Remove, now deprecated
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
	return s.nutter.Next()
}
