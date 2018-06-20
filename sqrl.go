package sqrl

import (
	"net/http"
)

// Nut returns a challenge
func Nut(r *http.Request) string {
	// 32 bits: user's connection IP address if secured, 0.0.0.0 if non-secured.
	// 32 bits: UNIX-time timestamp incrementing once per second.
	// 32 bits: up-counter incremented once for every SQRL link generated.
	// 31 bits: pseudo-random noise from system source.
	//  1  bit: flag bit to indicate source: QRcode or URL click

	return "TODO"
}

// User represents a users site specific public key.
type User string
