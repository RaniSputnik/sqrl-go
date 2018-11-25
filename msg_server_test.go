package sqrl_test

import (
	"testing"

	"github.com/RaniSputnik/sqrl-go"
	"github.com/stretchr/testify/assert"
)

func TestServerMsgEncode(t *testing.T) {
	t.Run("EncodesValidServerMsg", func(t *testing.T) {
		testCases := []struct {
			Input  sqrl.ServerMsg
			Expect string
		}{
			{
				Input: sqrl.ServerMsg{
					Ver: []string{sqrl.V1},
					Nut: "bRW-IegCUhGmcz9yvTtDKA",
					Tif: sqrl.TIF(5),
					Qry: "/sqrl?nut=bRW-IegCUhGmcz9yvTtDKA",
				},
				Expect: "dmVyPTENCm51dD1iUlctSWVnQ1VoR21jejl5dlR0REtBDQp0aWY9NQ0KcXJ5PS9zcXJsP251dD1iUlctSWVnQ1VoR21jejl5dlR0REtBDQo",
			},
		}

		for _, test := range testCases {
			got, err := test.Input.Encode()
			assert.NoError(t, err)
			assert.Equal(t, test.Expect, got)
		}
	})
}

func TestServerMsgParse(t *testing.T) {
	t.Run("ReturnsValidServerMsg", func(t *testing.T) {
		testCases := []struct {
			Input  string
			Expect sqrl.ServerMsg
		}{
			{
				Input: "dmVyPTENCm51dD1iUlctSWVnQ1VoR21jejl5dlR0REtBDQp0aWY9NQ0KcXJ5PS9zcXJsP251dD1iUlctSWVnQ1VoR21jejl5dlR0REtBDQo",
				Expect: sqrl.ServerMsg{
					Ver: []string{sqrl.V1},
					Nut: "bRW-IegCUhGmcz9yvTtDKA",
					Tif: 5,
					Qry: "/sqrl?nut=bRW-IegCUhGmcz9yvTtDKA",
				},
			},
			{
				Input: "dmVyPTENCm51dD1RTFlOd1N2TEZMZWd3RTlVMUZySG5BDQp0aWY9NA0KcXJ5PS9zcXJsP251dD1RTFlOd1N2TEZMZWd3RTlVMUZySG5BDQpzaW49MA0K",
				Expect: sqrl.ServerMsg{
					Ver: []string{sqrl.V1},
					Nut: "QLYNwSvLFLegwE9U1FrHnA",
					Tif: 4,
					Qry: "/sqrl?nut=QLYNwSvLFLegwE9U1FrHnA",
					// TODO: Sin: 0,
				},
			},
		}

		for _, test := range testCases {
			got, err := sqrl.ParseServer(test.Input)
			assert.NoError(t, err)
			if assert.NotNil(t, got) {
				assert.Equal(t, test.Expect, *got)
			}
		}
	})
}

func TestServerMsgIs(t *testing.T) {
	testCases := []struct {
		Input  sqrl.TIF
		Test   sqrl.TIF
		Expect bool
	}{
		{
			Input:  sqrl.TIFIPMatch,
			Test:   sqrl.TIFIPMatch,
			Expect: true,
		},
		{
			Input:  sqrl.TIFClientFailure | sqrl.TIFCommandFailed,
			Test:   sqrl.TIFCommandFailed,
			Expect: true,
		},
		{
			Input:  sqrl.TIFBadIDAssociation | sqrl.TIFIPMatch,
			Test:   sqrl.TIFIPMatch,
			Expect: true,
		},
		{
			Input:  sqrl.TIFFunctionNotSupported,
			Test:   sqrl.TIFIPMatch,
			Expect: false,
		},
		{
			Input:  sqrl.TIF(0),
			Test:   sqrl.TIFSQRLDisabled,
			Expect: false,
		},
	}

	for _, testCase := range testCases {
		serverMsg := sqrl.ServerMsg{Tif: testCase.Input}
		got := serverMsg.Is(testCase.Test)
		assert.Equal(t, testCase.Expect, got)
	}
}
