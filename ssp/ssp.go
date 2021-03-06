package ssp

import (
	"time"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/tomasen/realip"
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
	nutter := sqrl.NewNutter()

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

// WithClientEndpoint sets the endpoint that the client can
// use to post SQRL transactions to. This endpoint should
// be the path relative to the SQRL domain eg. /sqrl/cli.sqrl
//
// Defaults to /cli.sqrl if not set.
func (s *Server) WithClientEndpoint(url string) *Server {
	s.clientEndpoint = url
	return s
}

func (s *Server) Nut() sqrl.Nut {
	return s.nutter.Next()
}

// ClientIP is the function that is used to extract the client ip string
// from a given incomming http request. By default uses the FromRequest
// method from github.com/tomasen/realip, extracting the IP from either
// the X-Forwarded-For or X-Real-Ip headers, before falling back to remote addr.
var ClientIP = realip.FromRequest
