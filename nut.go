package sqrl

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"sync/atomic"

	"golang.org/x/crypto/blowfish"
)

// Nut is a base64, encrypted nonce that contains
// metadata about the request that it was derived from.
type Nut string

func (n Nut) String() string {
	return string(n)
}

// Nutter generates new nuts used to issue
// unique challenges to a SQRL client. It is
// also used to validate nuts that were
// previously issued.
type Nutter interface {
	Next() Nut
}

type blowfishNutter struct {
	cipher  *blowfish.Cipher
	counter uint32
}

// NewNutter creates a Nut generator
// with the given encryption key and a
// default nut expiry of 5 minutes.
// TODO: Key rotation
func NewNutter(key []byte) Nutter {
	cipher, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return &blowfishNutter{
		cipher: cipher,
	}
}

// Nut returns a challenge that should be returned to
// SQRL client for signing.
//
// The Nut (think nonce) is guaranteed to be unique
// and unpredictable to prevent replay attacks.
func (n *blowfishNutter) Next() Nut {
	nut := make([]byte, 8)

	// TODO combine this with a machine fingerprint
	count := atomic.AddUint32(&n.counter, 1)
	binary.LittleEndian.PutUint32(nut, count)

	noise := randBytes(4)
	nut[4] = noise[0]
	nut[5] = noise[1]
	nut[6] = noise[2]
	nut[7] = noise[3]

	encryptedNut := make([]byte, 8)
	n.cipher.Encrypt(encryptedNut, nut)
	return Nut(Base64.EncodeToString(encryptedNut))
}

func randBytes(length int) []byte {
	noise := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, noise); err != nil {
		// Nut generation does not currently return
		// an error as there is little recourse available
		// to a consumer.
		// It is probably safe to assume that a failure
		// to read random noise is a non-recoverable error.
		// This is an assumption that should be tested.
		panic(err.Error())
	}
	return noise
}
