package sqrl_test

import (
	"encoding/base64"
	"testing"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

func TestClientEncode(t *testing.T) {
	validIdk := sqrl.Identity("Vl4KVVRoG0C8v1VP0UEUNK2z_SYhNVYBXdoarhMljzQ")

	t.Run("FailsToEncodeInvalidClientMessages", func(t *testing.T) {
		cases := []struct {
			Name  string
			Input sqrl.ClientMsg
		}{
			{
				Name:  "Empty",
				Input: sqrl.ClientMsg{},
			},
			{
				Name: "MissingVersion",
				Input: sqrl.ClientMsg{
					Cmd: sqrl.CmdQuery,
					Idk: validIdk,
				},
			},
			{
				Name: "MissingCmd",
				Input: sqrl.ClientMsg{
					Ver: []string{sqrl.V1},
					Idk: validIdk,
				},
			},
			{
				Name: "MissingIdk",
				Input: sqrl.ClientMsg{
					Ver: []string{sqrl.V1},
					Cmd: sqrl.CmdQuery,
				},
			},
		}

		for _, test := range cases {
			_, err := test.Input.Encode()
			if err == nil {
				t.Errorf("%s: Expected encoding to fail, but got nil error", test.Name)
			}
		}
	})

	t.Run("EncodesValidClientMessagesCorrectly", func(t *testing.T) {
		cases := []struct {
			Name   string
			Input  sqrl.ClientMsg
			Expect string
		}{
			{
				Name: "Basic query",
				Input: sqrl.ClientMsg{
					Ver: []string{sqrl.V1},
					Cmd: sqrl.CmdQuery,
					Idk: validIdk,
				},
				Expect: "dmVyPTEKY21kPXF1ZXJ5Cmlkaz1WbDRLVlZSb0cwQzh2MVZQMFVFVU5LMnpfU1loTlZZQlhkb2FyaE1sanpRCg",
			},
		}

		for _, test := range cases {
			got, err := test.Input.Encode()
			if err != nil {
				t.Errorf("%s: Expected nil error, but got: '%v'", test.Name, err)
			}
			if got != test.Expect {
				t.Errorf("%s: Encoded value does not match\nExpected: '%s'\nGot:      '%s'",
					test.Name, test.Expect, got)
			}
		}
	})
}

func TestClientParse(t *testing.T) {

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
						// FIXME: Reinvestigate these tests
						//t.Errorf("Expected version at pos %d to be: %s, got: %s")
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
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(in))
}

func clientMessagesEqual(expected sqrl.ClientMsg, got sqrl.ClientMsg) bool {

	return true
}
