package ssp_test

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/RaniSputnik/sqrl-go/ssp"
)

const validNut = "rNRqu8olcWLAPaDvsL4b6owTVfryjzbre3hWHWnNTrK_hIS_KgIDFt2eBDc"

func TestHandlerNutIsReturned(t *testing.T) {
	s := httptest.NewServer(anyServer().Handler())
	res, err := http.Get(s.URL + "/nut.sqrl")

	// Assert no errors
	fatal(t, assert.NoError(t, err,
		"Expected no HTTP/connection error"))
	defer res.Body.Close()

	// Assert headers
	assert.Equal(t, http.StatusOK, res.StatusCode,
		"Expected successful status code")
	assert.True(t, strings.HasPrefix("application/x-www-form-urlencoded", res.Header.Get("Content-Type")),
		"Expected response to have Content-Type 'application/x-www-form-urlencoded'")

	values, err := parseNutResponse(res)
	fatal(t, assert.NoError(t, err, "Response error"))

	_, hasNut := values["nut"]
	assert.True(t, hasNut, "Missing nut parameter")
}

func parseNutResponse(res *http.Response) (url.Values, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %v", err)
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed parsing nut response body")
	}
	return values, nil
}
func TestHandlerNutIsUnique(t *testing.T) {
	s := httptest.NewServer(anyServer().Handler())
	endpoint := s.URL + "/nut.sqrl"

	results := make(map[string]bool)

	for i := 0; i < 10; i++ {
		res, err := http.Get(endpoint)
		fatal(t, assert.NoError(t, err,
			"Expected no HTTP/connection error"))
		defer res.Body.Close()

		values, err := parseNutResponse(res)
		fatal(t, assert.NoError(t, err, "Expected body to be decoded successfully"))

		seenBefore := results[values.Get("nut")]
		fatal(t, assert.False(t, seenBefore, "Duplicate nut returned: '%s'", values.Get("nut")))
		results[values.Get("nut")] = true
	}
}

func TestQRCodeIsReturned(t *testing.T) {
	s := httptest.NewServer(anyServer().Handler())
	res, err := http.Get(s.URL + "/qr.png?nut=" + validNut)

	// Assert no errors
	fatal(t, assert.NoError(t, err,
		"Expected no HTTP/connection error"))
	defer res.Body.Close()

	// Assert headers
	assert.Equal(t, http.StatusOK, res.StatusCode,
		"Expected successful status code")
	assert.Equal(t, "image/png", res.Header.Get("Content-Type"),
		"Expected response to have Content-Type 'image/png'")

	// Assert response body

	_, err = png.Decode(res.Body)
	fatal(t, assert.NoError(t, err,
		"Expected to decode the body as a PNG image successfully"))

	// TODO: We should compare images here to ensure the data was encoded successfully
}

func TestQRCodeIsReturnedAtSpecifiedSize(t *testing.T) {
	const givenSize = 64

	s := httptest.NewServer(anyServer().Handler())
	res, err := http.Get(fmt.Sprintf("%s/qr.png?nut=%s&size=%d", s.URL, validNut, givenSize))
	fatal(t, assert.NoError(t, err,
		"Expected no HTTP/connection error"))
	defer res.Body.Close()

	img, err := png.Decode(res.Body)
	fatal(t, assert.NoError(t, err,
		"Expected to decode the body as a PNG image successfully"))

	size := img.Bounds().Size()
	assert.Equal(t, givenSize, size.X, "Expected image width to match 'size' query parameter")
	assert.Equal(t, givenSize, size.Y, "Expected image width to match 'size' query parameter")
}

func fatal(t *testing.T, ok bool) {
	if !ok {
		t.FailNow()
	}
}

func anyServer() *ssp.Server {
	return ssp.Configure(make([]byte, 16), "http://example.com/auth/callback")
}

func anyTokenExchange() ssp.TokenExchange {
	return ssp.DefaultExchange(make([]byte, 16), time.Minute)
}
