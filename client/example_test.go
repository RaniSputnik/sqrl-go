package client_test

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/RaniSputnik/sqrl-go/client"
)

func TestExampleServer(t *testing.T) {
	httpClient := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/login", nil)
	fatal(t, err)

	resp, err := httpClient.Do(req)
	fatal(t, err)
	defer resp.Body.Close()

	challenge, err := extractSQRLChallenge(resp.Body, "localhost:8080")
	fatal(t, err)

	t.Logf("Challenge='%s'", challenge)
	if len(challenge) < 7 || challenge[:7] != "sqrl://" {
		t.Fatalf("Challenge is in unexpected format")
	}

	c := client.Client{UseInsecureConnection: true}
	expectErr(t, nil, c.Login(challenge))
}

func extractSQRLChallenge(body io.Reader, domain string) (string, error) {
	// Here, we search the body for a link with a value like the following
	// sqrl://www.grc.com/sqrl?nut=tx--bnqG7j4s-gEz1y4j8A
	// this is the SRQL challenge that we should answer in order to login.
	// TODO: don't read entire body, stop once we find the URL

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	bodyString := string(bytes)
	searchString := "sqrl://" + domain + "/sqrl?nut="
	left := strings.Index(bodyString, searchString)
	if left < 0 {
		return "", errors.New("SQRL link not found in body")
	}
	challenge := bodyString[left:]
	right := strings.Index(challenge, "\"")
	// Wierd bug in GRC's server
	if right2 := strings.Index(challenge, "<"); right2 < right {
		right = right2
	}
	challenge = challenge[:right]

	return challenge, nil
}

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected error: '<nil>', got: '%v'", err)
	}
}
