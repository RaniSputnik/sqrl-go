package sqrl

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/blowfish"
)

// Nutter generates new nuts used to issue
// unique challenges to a SQRL client. It is
// also used to validate nuts that were
// previously issued.
type Nutter struct {
	Expiry time.Duration

	key    []byte
	cipher *blowfish.Cipher
}

// NewNutter creates a Nut generator
// with the given encryption key and a
// default nut expiry of 5 minutes.
// TODO: Key rotation
func NewNutter(key []byte) *Nutter {
	cipher, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return &Nutter{
		key:    key,
		cipher: cipher,
		Expiry: time.Minute * 5,
	}
}

// NoClientID is used to represent a nut that will not
// perform any client identification check when validated.
const NoClientID = ""

var noClientIDBytes = make([]byte, 4)

var nuts uint32

// Nut is a base64, encrypted nonce that contains
// metadata about the request that it was derived from.
type Nut string

func (n Nut) String() string {
	return string(n)
}

// Nut returns a challenge that should be returned to
// SQRL client for signing.
//
// The Nut (think nonce) is guaranteed to be unique
// and unpredictable to prevent replay attacks.
func (n *Nutter) Nut(clientIdentifier string) Nut {
	nut := make([]byte, 8)

	// TODO combine this with a machine fingerprint
	count := atomic.AddUint32(&nuts, 1)
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

// Validate checks a nut returned by a client to ensure the nut
// is valid.
//
// The clients identifier (usually IP) is checked against the
// identifier encrypted in the nut to ensure the nut has been
// returned from the same machine it was originally sent to.
//
// Note: The client ID check will not be performed if the nut was
// created with NoClientID.
//
// The nut's expiry is also checked, to ensure there hasn't been
// a significant delay between nut issuing and nut return.
func (n *Nutter) Validate(returned Nut, clientIdentifier string) bool {
	decryptedNut, err := n.decryptNut(returned)
	if err != nil || len(decryptedNut) != 16 {
		return false // TODO: Do we need to expose this error?
	}

	originalIP := decryptedNut[:4]
	shouldCheckIP := bytes.Equal(originalIP, noClientIDBytes)
	if !shouldCheckIP {
		ip := nutClientIDBytes(clientIdentifier)
		if ipMatch := bytes.Equal(ip, originalIP); !ipMatch {
			return false
		}
	}

	timeSeconds := binary.BigEndian.Uint32(decryptedNut[4:8])
	t := time.Unix(int64(timeSeconds), 0)
	return time.Since(t) <= n.Expiry
}

func (n *Nutter) decryptNut(encrypted Nut) ([]byte, error) {
	decodedNut, err := Base64.DecodeString(string(encrypted))
	if err != nil {
		return nil, err
	}
	if len(decodedNut) != 8 {
		return nil, errors.New("invalid nut")
	}

	decryptedNut := make([]byte, 8)
	n.cipher.Decrypt(decryptedNut, decodedNut)
	// TODO: Verify the decryption was successful?
	return decodedNut, nil
}

func nutClientIDBytes(clientIdentifier string) []byte {
	if clientIdentifier == NoClientID {
		return noClientIDBytes
	}
	hashedClientID := md5.Sum([]byte(clientIdentifier))
	return hashedClientID[:4]
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
