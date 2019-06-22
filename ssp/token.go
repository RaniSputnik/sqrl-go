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

type token struct {
	UserId    string
	CreatedAt int64
}

func defaultNowFunc() time.Time {
	return time.Now()
}

// TODO: There is a lot of similarities here between sqrl.Server
// and the token generator - how could we share more of the logic
// between the two? Maybe managed aes is not required?
func NewTokenGenerator(key []byte) *TokenGenerator {
	aesgcm := genAesgcm(key)
	return &TokenGenerator{
		aesgcm:      aesgcm,
		tokenExpiry: time.Minute,
		NowFunc:     defaultNowFunc,
	}
}

type TokenGenerator struct {
	aesgcm      cipher.AEAD
	tokenExpiry time.Duration

	NowFunc func() time.Time
}

func (g *TokenGenerator) Token(userId string) string {
	if strings.Contains(userId, ",") {
		// Defensive against sneaky (and probably malicious) userIds
		panic("userId must not contain commas (,)")
	}
	t := g.encryptToken(token{
		UserId:    userId,
		CreatedAt: g.NowFunc().Unix(),
	})
	return string(t)
}

var (
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenFormatInvalid = errors.New("token format invalid")
)

func (g *TokenGenerator) ValidateToken(token string) (userId string, err error) {
	emptyUserId := ""
	t, err := g.decryptToken(token)
	if err != nil {
		return emptyUserId, ErrTokenFormatInvalid
	}
	issuedAt := time.Unix(t.CreatedAt, 0)
	if g.NowFunc().Sub(issuedAt) > g.tokenExpiry {
		return emptyUserId, ErrTokenExpired
	}
	// TODO: Any kind of validation needed of user id?
	return t.UserId, nil
}

var b64 = base64.URLEncoding.WithPadding(base64.NoPadding)

func (g *TokenGenerator) encryptToken(t token) string {
	payload := fmt.Sprintf("%s,%d", t.UserId, t.CreatedAt)
	nonce := randBytes(g.aesgcm.NonceSize())
	// TODO: Ensure payload is of suitable length
	encryptedToken := g.aesgcm.Seal(nil, nonce, []byte(payload), nil)
	encryptedTokenAndNonce := append(nonce, encryptedToken...)
	return b64.EncodeToString(encryptedTokenAndNonce)
}

func (g *TokenGenerator) decryptToken(t string) (token, error) {
	invalidToken := token{}
	encryptedTokenAndNonce, err := b64.DecodeString(t)
	if err != nil {
		return invalidToken, err
	}
	nonceSize := g.aesgcm.NonceSize()
	if len(encryptedTokenAndNonce) <= nonceSize {
		return invalidToken, errors.New("token length less than nonce size")
	}
	nonce := encryptedTokenAndNonce[:nonceSize]
	encryptedToken := encryptedTokenAndNonce[nonceSize:]

	decryptedToken, err := g.aesgcm.Open(nil, nonce, encryptedToken, nil)
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
