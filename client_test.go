package sqrl_test

import (
	"testing"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/stretchr/testify/assert"
)

const validIdk = sqrl.Identity("Vl4KVVRoG0C8v1VP0UEUNK2z_SYhNVYBXdoarhMljzQ")

func TestClientEncode(t *testing.T) {
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
				Input: "dmVyPTEKY21kPXF1ZXJ5Cmlkaz1WbDRLVlZSb0cwQzh2MVZQMFVFVU5LMnpfU1loTlZZQlhkb2FyaE1sanpRCg",
				Expected: sqrl.ClientMsg{
					Ver: []string{sqrl.V1},
					Cmd: sqrl.CmdQuery,
					Idk: validIdk,
				},
			},
		}

		for _, test := range cases {
			got, err := sqrl.ParseClient(test.Input)
			assert.NoError(t, err)
			if assert.NotNil(t, got) {
				assert.Equal(t, test.Expected, *got)
			}
		}
	})
}
