package client_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/RaniSputnik/sqrl-go/client"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:61.0) Gecko/20100101 Firefox/61.0"

//const cookies = "pag=zngakt1fkp0es; pcss=3wmopr5i4kblq; pico=3hsz2fslhls52; tpag=zngakt1fkp0es; sqrluser=2oqdgwdprpi5o; tcss=3wmopr5i4kblq; tico=3hsz2fslhls52"
const cookies = "ppag=sq5ln315sn2uw; pcss=iknrukwyhaf5u; pico=3thoa0rforkor; sqrluser=f2uheyqefkf4k; tpag=sq5ln315sn2uw; tcss=iknrukwyhaf5u; tico=3thoa0rforkor"

// TOOD gorace - keep getting 'Post https://www.grc.com/sqrl?nut=XXX: EOF' error
// But when I debug the error does not occur. Presume race condition.
func TestLiveServer(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://www.grc.com/sqrl/diag.htm", nil)
	fatal(t, err)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Cookie", cookies)

	resp, err := http.DefaultClient.Do(req)
	fatal(t, err)
	defer resp.Body.Close()

	challenge, err := extractSQRLChallenge(resp.Body)
	fatal(t, err)

	t.Logf("Challenge='%s'", challenge)
	if len(challenge) < 7 || challenge[:7] != "sqrl://" {
		t.Fatalf("Challenge is in unexpected format")
	}

	expectErr(t, nil, client.Login(challenge))
}

func extractSQRLChallenge(body io.Reader) (string, error) {
	// Here, we search the body for a link with a value like the following
	// sqrl://www.grc.com/sqrl?nut=tx--bnqG7j4s-gEz1y4j8A
	// this is the SRQL challenge that we should answer in order to login.
	// TODO: don't read entire body, stop once we find the URL

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	bodyString := string(bytes)
	searchString := "sqrl://www.grc.com/sqrl?nut="
	left := strings.Index(bodyString, searchString)
	challenge := bodyString[left:]
	right := strings.Index(challenge, "\"")
	challenge = challenge[:right]

	return challenge, nil
}

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected error: '<nil>', got: '%v'", err)
	}
}
