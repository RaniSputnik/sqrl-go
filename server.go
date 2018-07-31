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
	Nut string // TODO: type Nut
	Tif TIF
	Qry string

	// TODO: url - for any command other than query
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
		"nut=" + m.Nut,
		"tif=" + strconv.Itoa(int(m.Tif)),
		"qry=" + m.Qry,
		"", // Must end with a final newline
	}
	return Base64.EncodeToString([]byte(strings.Join(vals, "\r\n"))), nil
}

// ParseServer decodes the base64 encoded server
// parameter into the component parts.
func ParseServer(raw string) (*ServerMsg, error) {
	vals, err := parseMsg(raw)
	if len(vals) < 4 {
		return nil, errors.New("missing one or more required parameters (ver,nut,tif,qry)")
	}

	tifstr := vals["tif"]
	tif, err := strconv.Atoi(tifstr)
	if err != nil {
		return nil, fmt.Errorf("required value 'tif' is invalid: '%s'", tifstr)
	}

	// TODO: Check supported version before parsing
	return &ServerMsg{
		Ver: parseVer(vals["ver"]),
		Nut: vals["nut"],
		Tif: TIF(tif),
		Qry: vals["qry"],
	}, nil
}
