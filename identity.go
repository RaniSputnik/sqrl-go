package sqrl

import (
	"golang.org/x/crypto/ed25519"
)

// Identity represents a users site specific public key,
// base64 encoded for transmission.
type Identity string

// Signature is a base64 signature sent by the
// client. The signature can be verified using the
// corresponding identity.
type Signature string

// Verify determines if the given signature is valid.
func (s Signature) Verify(id Identity, payload string) bool {
	publicKey, err := Base64.DecodeString(string(id))
	if err != nil {
		return false
	}
	signature, err := Base64.DecodeString(string(s))
	if err != nil {
		return false
	}
	return ed25519.Verify(publicKey, []byte(payload), signature)
}
