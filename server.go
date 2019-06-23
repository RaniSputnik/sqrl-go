package sqrl

import (
	"crypto/aes"
	"crypto/cipher"
	"time"
)

// TODO: Change this to "nut generator"
// It's no longer a server, it should have
// a more focussed responsibility

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
