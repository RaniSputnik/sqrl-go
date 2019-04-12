package sqrl

import (
	"errors"
	"strings"
)

// ClientMsg is used to represent the values
// sent by the client to the server.
//
// Ver are the versions of SQRL the client supports.
// Cmd is the command type for this request.
type ClientMsg struct {
	Ver []string
	Cmd Cmd
	Idk Identity

	Opt []Opt
}

// Encode writes the client message to a string
// ready for transmission to a SQRL server.
func (m *ClientMsg) Encode() (string, error) {
	if len(m.Ver) == 0 || m.Cmd == "" || m.Idk == "" {
		return "", errors.New("incomplete client message")
	}

	vals := []string{
		"ver=" + strings.Join(m.Ver, ","),
		"cmd=" + string(m.Cmd),
		"idk=" + string(m.Idk),
	}
	if len(m.Opt) > 0 {
		vals = append(vals, "opt="+encodeOptions(m.Opt))
	}
	vals = append(vals, "") // Must end with a final newline
	return Base64.EncodeToString([]byte(strings.Join(vals, "\r\n"))), nil
}

func encodeOptions(opts []Opt) string {
	res := ""
	lastOpt := len(opts) - 1
	for i, opt := range opts {
		res += string(opt)
		if i != lastOpt {
			res += "~"
		}
	}
	return res
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

// ParseClient decodes a client message from the given string
func ParseClient(raw string) (*ClientMsg, error) {
	vals, err := parseMsg(raw)
	if err != nil {
		return nil, err
	}

	if vals["ver"] == "" || vals["cmd"] == "" || vals["idk"] == "" {
		return nil, errors.New("missing required parameter")
	}

	ver, err := parseVer(vals["ver"])
	if err != nil {
		return nil, err
	}

	return &ClientMsg{
		Ver: ver,
		Cmd: Cmd(vals["cmd"]),
		Idk: Identity(vals["idk"]),
		Opt: parseOpts(vals["opt"]),
	}, nil
}

func parseOpts(input string) []Opt {
	if input == "" {
		return []Opt{}
	}

	// TODO: Remove strings.Split and reduce allocations
	inputs := strings.Split(input, "~")
	opts := make([]Opt, len(inputs))
	for i, optString := range inputs {
		opts[i] = Opt(optString)
	}

	// TODO: Do we want to check for validity of options?
	// Or do we simply ignore options we don't understand?

	return opts
}
