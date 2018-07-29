package sqrl

import (
	"errors"
	"fmt"
	"strconv"
)

// ServerMsg is used to represent the values
// sent from the server to the client.
type ServerMsg struct {
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
