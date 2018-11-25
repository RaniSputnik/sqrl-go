package sqrlhttp

import (
	"fmt"
	"net/http"

	sqrl "github.com/RaniSputnik/sqrl-go"
	qrcode "github.com/skip2/go-qrcode"
)

func GenerateChallenge(server *sqrl.Server, r *http.Request, domain string) (url string, qr []byte) {
	nonce := server.Nut(clientID(r))
	loginURL := fmt.Sprintf("sqrl://%s/sqrl?nut=%s", domain, nonce)
	qrCode, _ := qrcode.Encode(loginURL, qrcode.Medium, 256) // TODO: What to do if QR Code generation fails?
	return loginURL, qrCode
}
