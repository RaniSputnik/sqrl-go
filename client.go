package sqrl

import (
	"encoding/base64"
	"errors"
)

// ClientMsg is used to represent the values
// sent by the client to the server.
//
// Ver are the versions of SQRL the client supports.
// Cmd is the command type for this request.
type ClientMsg struct {
	Ver []string
	Cmd Cmd

	Opt []Opt
}

// HasOpt returns true/false whether the given option was provided.
func (m ClientMsg) HasOpt(query Opt) bool {
	if m.Opt == nil {
		return false
	}
	for _, o := range m.Opt {
		if o == query {
			return true
		}
	}
	return false
}

func ParseClient(raw string) (ClientMsg, error) {
	_, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return ClientMsg{}, err
	}

	return ClientMsg{}, errors.New("Not implemented")
}
