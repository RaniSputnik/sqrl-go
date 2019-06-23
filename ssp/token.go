package ssp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type Token string

type TokenGenerator interface {
	Token(userId string) Token
}

type TokenValidator interface {
	Validate(token Token) (userId string, err error)
}

type TokenExchange interface {
	TokenGenerator
	TokenValidator
}

type token struct {
	UserId    string
	CreatedAt int64
}

// TODO: There is a lot of similarities here between sqrl.Server
// and the token generator - how could we share more of the logic
// between the two? Maybe managed aes is not required?
func DefaultExchange(key []byte, expiry time.Duration) TokenExchange {
	aesgcm := genAesgcm(key)
	return &defaultExchange{
		aesgcm:      aesgcm,
		tokenExpiry: expiry,
	}
}

type defaultExchange struct {
	aesgcm      cipher.AEAD
	tokenExpiry time.Duration
}

func (g *defaultExchange) Token(userId string) Token {
	if strings.Contains(userId, ",") {
		// Defensive against sneaky (and probably malicious) userIds
		panic("userId must not contain commas (,)")
	}
	t := g.encryptToken(token{
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
	})
	return Token(t)
}

var (
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenFormatInvalid = errors.New("token format invalid")
)

func (e *defaultExchange) Validate(token Token) (userId string, err error) {
	t, err := e.decryptToken(token)
	if err != nil {
		return "", ErrTokenFormatInvalid
	}
	issuedAt := time.Unix(t.CreatedAt, 0)
	if time.Since(issuedAt) > e.tokenExpiry {
		return "", ErrTokenExpired
	}
	// TODO: Any kind of validation needed of user id?
	return t.UserId, nil
}

var b64 = base64.URLEncoding.WithPadding(base64.NoPadding)

func (e *defaultExchange) encryptToken(t token) Token {
	payload := fmt.Sprintf("%s,%d", t.UserId, t.CreatedAt)
	nonce := randBytes(e.aesgcm.NonceSize())
	// TODO: Ensure payload is of suitable length
	encryptedToken := e.aesgcm.Seal(nil, nonce, []byte(payload), nil)
	encryptedTokenAndNonce := append(nonce, encryptedToken...)
	return Token(b64.EncodeToString(encryptedTokenAndNonce))
}

func (e *defaultExchange) decryptToken(t Token) (token, error) {
	invalidToken := token{}
	encryptedTokenAndNonce, err := b64.DecodeString(string(t))
	if err != nil {
		return invalidToken, err
	}
	nonceSize := e.aesgcm.NonceSize()
	if len(encryptedTokenAndNonce) <= nonceSize {
		return invalidToken, errors.New("token length less than nonce size")
	}
	nonce := encryptedTokenAndNonce[:nonceSize]
	encryptedToken := encryptedTokenAndNonce[nonceSize:]

	decryptedToken, err := e.aesgcm.Open(nil, nonce, encryptedToken, nil)
	if err != nil {
		return invalidToken, err
	}

	tokenParts := strings.Split(string(decryptedToken), ",")
	if len(tokenParts) != 2 {
		return invalidToken, errors.New("token should have exactly one comma")
	}
	userId := tokenParts[0]
	createdAt, err := strconv.ParseInt(tokenParts[1], 10, 64)
	if err != nil {
		return invalidToken, err
	}

	return token{
		UserId:    userId,
		CreatedAt: createdAt,
	}, nil
}

func randBytes(length int) []byte {
	noise := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, noise); err != nil {
		// Token generation does not currently return
		// an error as there is little recourse available
		// to a consumer.
		// It is probably safe to assume that a failure
		// to read random noise is a non-recoverable error.
		// This is an assumption that should be tested.
		panic(err.Error())
	}
	return noise
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
