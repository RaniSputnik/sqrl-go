package sqrl_test

import (
	"fmt"
	"testing"

	sqrl "github.com/RaniSputnik/sqrl-go"
	"github.com/stretchr/testify/assert"
)

const validIdk = sqrl.Identity("Vl4KVVRoG0C8v1VP0UEUNK2z_SYhNVYBXdoarhMljzQ")

func TestClientMsgEncode(t *testing.T) {
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
				Expect: "dmVyPTENCmNtZD1xdWVyeQ0KaWRrPVZsNEtWVlJvRzBDOHYxVlAwVUVVTksyel9TWWhOVllCWGRvYXJoTWxqelENCg",
			},
			{
				Name: "Query with single option",
				Input: sqrl.ClientMsg{
					Ver: []string{sqrl.V1},
					Cmd: sqrl.CmdQuery,
					Idk: validIdk,
					Opt: []sqrl.Opt{sqrl.OptCPS},
				},
				Expect: "dmVyPTENCmNtZD1xdWVyeQ0KaWRrPVZsNEtWVlJvRzBDOHYxVlAwVUVVTksyel9TWWhOVllCWGRvYXJoTWxqelENCm9wdD1jcHMNCg",
			},
			{
				Name: "Query with two options",
				Input: sqrl.ClientMsg{
					Ver: []string{sqrl.V1},
					Cmd: sqrl.CmdQuery,
					Idk: validIdk,
					Opt: []sqrl.Opt{sqrl.OptCPS, sqrl.OptSQRLOnly},
				},
				Expect: "dmVyPTENCmNtZD1xdWVyeQ0KaWRrPVZsNEtWVlJvRzBDOHYxVlAwVUVVTksyel9TWWhOVllCWGRvYXJoTWxqelENCm9wdD1jcHN-c3FybG9ubHkNCg",
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

func TestClientMsgParse(t *testing.T) {
	t.Run("ReturnsErrorWhenClientStringInvalid", func(t *testing.T) {
		validIdk := "PO2ib4BeITiKHTOGW37Mv03dES29DfhJPl6bq5JijoA"

		cases := []struct {
			Name  string
			Input string
		}{
			{"Empty", ""},
			{"NotBase64", "notbase64!!@!@£$"},
			{"OnlyWhitespace", "         "},
			{"DuplicateFields", sqrl.Base64.EncodeToString([]byte("ver=1\nver=1\ncmd=query\nidk=" + validIdk))},
			{"MissingIdkField", sqrl.Base64.EncodeToString([]byte("ver=1\ncmd=query"))},
			{"MissingCmdField", sqrl.Base64.EncodeToString([]byte("ver=1\nidk=" + validIdk))},
			{"MissingVerField", sqrl.Base64.EncodeToString([]byte("cmd=query\nidk=" + validIdk))},
			// The Web extension does not lead with the version information
			// TODO: Is this a bug in the extension? Or should we relax this constraint?
			// https://github.com/RaniSputnik/sqrl-go/issues/12
			// {"VerNotFirst", sqrl.Base64.EncodeToString([]byte("cmd=query\nver=1\nidk=" + validIdk))},
		}

		for _, testCase := range cases {
			t.Run(testCase.Name, func(t *testing.T) {
				_, err := sqrl.ParseClient(testCase.Input)
				assert.Error(t, err)
			})
		}
	})

	t.Run("ReturnsValidClient", func(t *testing.T) {
		cases := []struct {
			Name     string
			Input    string
			Expected sqrl.ClientMsg
		}{
			{
				Name:  "Basic query",
				Input: "dmVyPTEKY21kPXF1ZXJ5Cmlkaz1WbDRLVlZSb0cwQzh2MVZQMFVFVU5LMnpfU1loTlZZQlhkb2FyaE1sanpRCg",
				Expected: sqrl.ClientMsg{
					Ver: []string{sqrl.V1},
					Cmd: sqrl.CmdQuery,
					Idk: validIdk,
					Opt: []sqrl.Opt{},
				},
			},
			{
				Name:  "Basic query, alternative newline characters",
				Input: "dmVyPTENCmNtZD1xdWVyeQ0KaWRrPVZsNEtWVlJvRzBDOHYxVlAwVUVVTksyel9TWWhOVllCWGRvYXJoTWxqelENCg",
				Expected: sqrl.ClientMsg{
					Ver: []string{sqrl.V1},
					Cmd: sqrl.CmdQuery,
					Idk: validIdk,
					Opt: []sqrl.Opt{},
				},
			},
			{
				Name:  "Query with single option",
				Input: "dmVyPTENCmNtZD1xdWVyeQ0KaWRrPVZsNEtWVlJvRzBDOHYxVlAwVUVVTksyel9TWWhOVllCWGRvYXJoTWxqelENCm9wdD1jcHMNCg",
				Expected: sqrl.ClientMsg{
					Ver: []string{sqrl.V1},
					Cmd: sqrl.CmdQuery,
					Idk: validIdk,
					Opt: []sqrl.Opt{sqrl.OptCPS},
				},
			},
			{
				Name:  "Query with two options",
				Input: "dmVyPTENCmNtZD1xdWVyeQ0KaWRrPVZsNEtWVlJvRzBDOHYxVlAwVUVVTksyel9TWWhOVllCWGRvYXJoTWxqelENCm9wdD1jcHN-c3FybG9ubHkNCg",
				Expected: sqrl.ClientMsg{
					Ver: []string{sqrl.V1},
					Cmd: sqrl.CmdQuery,
					Idk: validIdk,
					Opt: []sqrl.Opt{sqrl.OptCPS, sqrl.OptSQRLOnly},
				},
			},
		}

		for _, test := range cases {
			got, err := sqrl.ParseClient(test.Input)
			assert.NoError(t, err, fmt.Sprintf("'%s' failed, returned an error", test.Name))
			if assert.NotNil(t, got, fmt.Sprintf("'%s' failed, result was nil", test.Name)) {
				assert.Equal(t, test.Expected, *got,
					fmt.Sprintf("'%s' failed, result did not match expected value", test.Name))
			}
		}
	})
}
