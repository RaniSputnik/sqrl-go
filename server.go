package sqrl

import (
	"crypto/aes"
	"crypto/cipher"
	"strings"
	"time"
)

// Server is a SQRL compliant server configured
// with nut encryption keys and expiry times.
type Server struct {
	key       []byte
	aesgcm    cipher.AEAD
	nutExpiry time.Duration

	redirectURL    string
	clientEndpoint string
}

// Key returns the AES key used
// by the server for nut generation
//
// This is currently used by the
// SSP package to share a single key
// between servers. Can we do better
// than this?
func (s *Server) Key() []byte {
	return s.key
}

// Configure creates a new SQRL server
// with the given encryption key and a
// default nut expiry of 5 minutes.
// TODO: Key rotation
func Configure(key []byte) *Server {
	aesgcm := genAesgcm(key)
	return &Server{
		key:            key,
		aesgcm:         aesgcm,
		nutExpiry:      time.Minute * 5,
		redirectURL:    "/",
		clientEndpoint: "/cli.sqrl",
	}
}

func genAesgcm(key []byte) cipher.AEAD {
	padKeyIfRequired(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	return aesgcm
}

func padKeyIfRequired(key []byte) {
	// TODO: Ensure key is either 16, 24 or 32 bits
}

// WithNutExpiry sets the window of time within which
// a nut is considered to be valid.
func (s *Server) WithNutExpiry(d time.Duration) *Server {
	s.nutExpiry = d
	return s
}

// WithRedirectURL sets the endpoint that the SQRL server
// will redirect to when authentication is successful.
//
// When redirecting the SQRL server will add a ? then the
// authentication token that can be used to retrieve the
// session eg. If the URL https://example.com is provided
// tokens will be returned to https://example.com?12345678
func (s *Server) WithRedirectURL(url string) *Server {
	i := strings.LastIndexByte(url, '?')
	if i > -1 {
		s.redirectURL = url[:i]
	} else {
		s.redirectURL = url
	}
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

// TODO: Do we need to expose these getters?
// Complicates the interface a little, maybe would
// be better as a simple struct?

func (s *Server) NutExpiry() time.Duration {
	return s.nutExpiry
}

func (s *Server) RedirectURL() string {
	return s.redirectURL
}

func (s *Server) ClientEndpoint() string {
	return s.clientEndpoint
}
