// +build ignore

package client_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/RaniSputnik/sqrl-go/client"
)

const cookies = "ppag=sq5ln315sn2uw; pcss=iknrukwyhaf5u; pico=3thoa0rforkor; sqrluser=f2uheyqefkf4k; tpag=sq5ln315sn2uw; tcss=iknrukwyhaf5u; tico=3thoa0rforkor"

func TestLiveServer(t *testing.T) {
	httpClient := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest(http.MethodGet, "https://www.grc.com/sqrl/diag.htm", nil)
	fatal(t, err)
	req.Header.Add("Cookie", cookies)

	resp, err := httpClient.Do(req)
	fatal(t, err)
	defer resp.Body.Close()

	challenge, err := extractSQRLChallenge(resp.Body)
	fatal(t, err)

	t.Logf("Challenge='%s'", challenge)
	if len(challenge) < 7 || challenge[:7] != "sqrl://" {
		t.Fatalf("Challenge is in unexpected format")
	}

	waitForGRCWaitLimit()
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
	// Wierd bug in GRC's server
	if right2 := strings.Index(challenge, "<"); right2 < right {
		right = right2
	}
	challenge = challenge[:right]

	return challenge, nil
}

func waitForGRCWaitLimit() {
	time.Sleep(1 * time.Second)
}

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected error: '<nil>', got: '%v'", err)
	}
}
