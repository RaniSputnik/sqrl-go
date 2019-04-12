package sqrl

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ServerMsg is used to represent the values
// sent from the server to the client.
type ServerMsg struct {
	// TODO: should we remove Ver? Do package
	// consumers actually need to be able to
	// set this?
	Ver []string
	Nut Nut
	Tif TIF
	Qry string
	URL string

	// TODO: sin - Secret index
	// TODO: suk - server unlock key
	// TODO: ask - message text to display to user
	// TODO: can - cancellation redirection URL

	// TODO: additional parameters
}

// Encode writes the server message to a string
// ready for transmission to the client.
func (m *ServerMsg) Encode() (string, error) {
	vals := []string{
		"ver=" + strings.Join(m.Ver, ","),
		"nut=" + string(m.Nut),
		"tif=" + strconv.Itoa(int(m.Tif)),
		"qry=" + m.Qry,
	}
	if m.URL != "" {
		vals = append(vals, "url="+m.URL)
	}
	vals = append(vals, "") // Must end with a final newline
	return Base64.EncodeToString([]byte(strings.Join(vals, "\r\n"))), nil
}

// Set adds the given transaction information flag
// to the server message.
func (m *ServerMsg) Set(flag TIF) {
	m.Tif |= flag
}

// Unset removes the given transaction information
// flag from the server message.
func (m *ServerMsg) Unset(flag TIF) {
	m.Tif &^= flag
}

// Is returns whether or not the server message
// includes the given transaction information flag.
func (m *ServerMsg) Is(flag TIF) bool {
	return m.Tif&flag != 0
}

// ParseServer decodes the base64 encoded server
// parameter into the component parts.
func ParseServer(raw string) (*ServerMsg, error) {
	vals, err := parseMsg(raw)
	if err != nil {
		return nil, err
	}

	if len(vals) < 4 {
		return nil, errors.New("missing one or more required parameters (ver,nut,tif,qry)")
	}

	ver, err := parseVer(vals["ver"])
	if err != nil {
		return nil, err
	}

	tifstr := vals["tif"]
	tif, err := strconv.Atoi(tifstr)
	if err != nil {
		return nil, fmt.Errorf("required value 'tif' is invalid: '%s'", tifstr)
	}

	// TODO: Ensure nut can be decoded correctly
	nut := Nut(vals["nut"])

	// TODO: Check supported version before parsing
	return &ServerMsg{
		Ver: ver,
		Nut: nut,
		Tif: TIF(tif),
		Qry: vals["qry"],
		URL: vals["url"],
	}, nil
}
