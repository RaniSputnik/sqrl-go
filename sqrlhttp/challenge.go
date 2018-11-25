package sqrlhttp

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	sqrl "github.com/RaniSputnik/sqrl-go"
	qrcode "github.com/skip2/go-qrcode"
)

func GenerateChallenge(server *sqrl.Server, r *http.Request, domain string) template.HTML {
	nonce := server.Nut(clientID(r))
	loginURL := fmt.Sprintf("sqrl://%s/sqrl?nut=%s", domain, nonce)
	qrCode, _ := qrcode.Encode(loginURL, qrcode.Medium, 256) // TODO: What to do if QR Code generation fails?
	qrCodeSrc := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(qrCode))
	const linkFormat = `<a href="%s" target="_blank"><img src="%s" alt="SQRL Login" /></a>`
	const jsFormat = `<script type="text/javascript">` + syncJS + `</script>`
	fragment := fmt.Sprintf(linkFormat+jsFormat, loginURL, qrCodeSrc)
	return template.HTML(fragment)
}
