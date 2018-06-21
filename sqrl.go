package sqrl

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

var nuts uint32

// Nut returns a challenge
func Nut(r *http.Request) string {
	//  32 bits: user's connection IP address if secured, 0.0.0.0 if non-secured.
	//  32 bits: UNIX-time timestamp incrementing once per second.
	//  32 bits: up-counter incremented once for every SQRL link generated.
	//  31 bits: pseudo-random noise from system source.
	//   1  bit: flag bit to indicate source: QRcode or URL click
	// ---------
	// 128 bits: AES encryption block size

	nut := make([]byte, 16)

	//  32 bits: user's connection IP address if secured, 0.0.0.0 if non-secured.
	// TODO X-Forwarded-For
	// TODO only parse IP Address if it is a secure connection
	ip := parseIP(r.RemoteAddr)
	if !ip.IsLoopback() {
		ip = ip.To4()
		if len(ip) == net.IPv4len {
			nut[0] = ip[0]
			nut[1] = ip[1]
			nut[2] = ip[2]
			nut[3] = ip[3]
		}
	}

	// UNIX-time timestamp incrementing once per second.
	time := uint32(time.Now().Unix())
	nut[4] = byte(time)
	nut[5] = byte(time >> 8)
	nut[6] = byte(time >> 16)
	nut[7] = byte(time >> 24)

	//  32 bits: up-counter incremented once for every SQRL link generated.
	// TODO combine this with a machine fingerprint
	count := atomic.AddUint32(&nuts, 1)
	nut[8] = byte(count >> 24)
	nut[9] = byte(count >> 16)
	nut[10] = byte(count >> 8)
	nut[11] = byte(count)

	//  31 bits: pseudo-random noise from system source.
	noise := make([]byte, 4)
	if _, err := io.ReadFull(rand.Reader, noise); err != nil {
		panic(err.Error())
	}
	nut[12] = noise[0]
	nut[13] = noise[1]
	nut[14] = noise[2]
	nut[15] = noise[3]

	//   1  bit: flag bit to indicate source: QRcode or URL click
	// TODO

	// TODO gen key
	// TODO key rotation
	key := make([]byte, 16)
	block, err := aes.NewCipher(key)

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	// TODO should we base64 encode here instead?
	return hex.EncodeToString(aesgcm.Seal(nil, nonce, nut, nil))
}

func parseIP(remoteAddr string) net.IP {
	// TODO this func is rubbish, clean it up
	res := remoteAddr
	if remoteAddr[0] == '[' {
		i := strings.LastIndex(remoteAddr, "]")
		if i > 0 {
			res = remoteAddr[1:i]
		} else {
			res = remoteAddr[1:]
		}
	}
	if pci := strings.IndexRune(res, '%'); pci > -1 {
		res = res[:pci]
	}
	log.Print(res)
	// TODO strip port number
	return net.ParseIP(res)
}

// User represents a users site specific public key.
type User string
