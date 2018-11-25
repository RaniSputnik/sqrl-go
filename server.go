package sqrl

import (
	"crypto/aes"
	"crypto/cipher"
	"time"
)

// Server is a SQRL compliant server configured
// with nut encryption keys and expiry times.
type Server struct {
	aesgcm    cipher.AEAD
	nutExpiry time.Duration
}

// Configure creates a new SQRL server
// with the given encryption key and a
// default nut expiry of 5 minutes.
// TODO: Key rotation
func Configure(key []byte) *Server {
	aesgcm := genAesgcm(key)
	return &Server{aesgcm, time.Minute * 5}
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
	// TODO: Mutex?
	s.nutExpiry = d
	return s
}
