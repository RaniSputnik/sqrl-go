package sqrl

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Server is a SQRL compliant server configured
// with nut encryption keys and expiry times.
type Server struct {
	key       []byte
	nutExpiry time.Duration
}

// Configure creates a new SQRL server
// with the given encryption key and a
// default nut expiry of 5 minutes.
// TODO: Key rotation
func Configure(key []byte) *Server {
	padKeyIfRequired(key)
	return &Server{key, time.Minute * 5}
}

func padKeyIfRequired(key []byte) {
}

// WithNutExpiry sets the window of time within which
// a nut is considered to be valid.
func (s *Server) WithNutExpiry(d time.Duration) *Server {
	// TODO: Mutex?
	s.nutExpiry = d
	return s
}

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
		"nut=" + string(m.Nut),
		"tif=" + strconv.Itoa(int(m.Tif)),
		"qry=" + m.Qry,
		"", // Must end with a final newline
	}
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
	}, nil
}
