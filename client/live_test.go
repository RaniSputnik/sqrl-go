// +build ignore

package client_test

import (
	"net/http"
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

	challenge, err := extractSQRLChallenge(resp.Body, "www.grc.com")
	fatal(t, err)

	t.Logf("Challenge='%s'", challenge)
	if len(challenge) < 7 || challenge[:7] != "sqrl://" {
		t.Fatalf("Challenge is in unexpected format")
	}

	waitForGRCWaitLimit()
	expectErr(t, nil, client.Login(challenge))
}

func waitForGRCWaitLimit() {
	time.Sleep(1 * time.Second)
}
