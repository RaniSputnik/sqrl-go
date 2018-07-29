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
	Ver []string
	Nut string // TODO: type Nut
	Tif int
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
	bytes, err := Base64.DecodeString(raw)
	if err != nil {
		return nil, err
	}

	form := strings.Split(string(bytes), "\r\n") // TODO: Move newline char to const
	if len(form) < 4 {
		return nil, errors.New("missing one or more required parameters (ver,nut,tif,qry)")
	}

	vals := map[string]string{}
	for _, keyval := range form {
		if keyval == "" {
			continue
		}

		pair := strings.SplitN(keyval, "=", 2)
		if len(pair) < 2 {
			return nil, fmt.Errorf("invalid value '%s', should be in the form: key=value\\r\\n", keyval)
		}
		vals[pair[0]] = pair[1]
	}

	tifstr := vals["tif"]
	tif, err := strconv.Atoi(tifstr)
	if err != nil {
		return nil, fmt.Errorf("required value 'tif' is invalid: '%s'", tifstr)
	}

	// TODO: Check supported version before parsing
	return &ServerMsg{
		Ver: strings.Split(vals["ver"], ","),
		Nut: vals["nut"],
		Tif: tif,
		Qry: vals["qry"],
	}, nil
}
