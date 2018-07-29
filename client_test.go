package sqrl_test

import (
	"encoding/base64"
	"testing"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

func TestParseClient(t *testing.T) {

	t.Run("ReturnsErrorWhenClientStringInvalid", func(t *testing.T) {

	})

	t.Run("ReturnsValidClient", func(t *testing.T) {
		cases := []struct {
			Input    string
			Expected sqrl.ClientMsg
		}{
			{
				Input: "ver=1\ncmd=query",
				Expected: sqrl.ClientMsg{
					Ver: []string{sqrl.V1},
					Cmd: sqrl.CmdQuery,
				},
			},
		}

		for _, test := range cases {
			expect := test.Expected

			got, err := sqrl.ParseClient(b64(test.Input))
			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if expected := len(expect.Ver); expected != len(got.Ver) {
				t.Errorf("Expected: %d versions, got: %d", expected, len(got.Ver))
			} else {
				for i := 0; i < len(expect.Ver); i++ {
					if expect.Ver[i] != got.Ver[i] {
						t.Errorf("Expected version at pos %d to be: %s, got: %s")
					}
				}
			}

			if expect.Cmd != got.Cmd {
				t.Errorf("Expected cmd: %s, got: %s", expect.Cmd, got.Cmd)
			}
		}
	})
}

func b64(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}

func clientMessagesEqual(expected sqrl.ClientMsg, got sqrl.ClientMsg) bool {

	return true
}
