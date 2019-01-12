package ssp

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	sqrl "github.com/RaniSputnik/sqrl-go"
	qrcode "github.com/skip2/go-qrcode"
)

func clientID(r *http.Request) string {
	// TODO: X-Forwarded-For
	// TODO: Include user agent if available
	return r.RemoteAddr
}

func requestDomain(r *http.Request) string {
	// TODO: Do we need to do anything special here for proxies?
	return r.Host
}

// GenerateChallenge creates a login URL that a SQRL client can
// use to perform login. It additionally generates a HTML fragment
// that should be rendered into the page and presented to the user.
//
// The fragment contains a QR code that is wrapped in a sqrl:// link
// and a script that will poll the server for authentication status.
// The script will redirect the user once login has been completed.
func GenerateChallenge(server *sqrl.Server, r *http.Request) (string, template.HTML) {
	// TODO: How might we remove QR code generation from here and keep the API super simple
	// TODO: How might we include the sync endpoint handling in the ssp package?
	nonce := server.Nut(clientID(r))
	loginURL := fmt.Sprintf("sqrl://%s/sqrl?nut=%s", requestDomain(r), nonce)
	qrCode, _ := qrcode.Encode(loginURL, qrcode.Medium, 256) // TODO: What to do if QR Code generation fails?
	qrCodeSrc := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(qrCode))
	const linkFormat = `<a href="%s" target="_blank"><img src="%s" alt="SQRL Login" /></a>`
	const jsFormat = `<script type="text/javascript">` + syncJS + `</script>`
	fragment := fmt.Sprintf(linkFormat+jsFormat, loginURL, qrCodeSrc)
	return loginURL, template.HTML(fragment)
}
