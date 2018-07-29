package sqrl

// User represents a users site specific public key.
type User string

const (
	// V1 is the version 1 string for the SQRL protocol.
	V1 = "1"

	// Scheme is sqrl: used in SQRL URI's
	Scheme = "sqrl"

	// TODO what to do about qrl://
	// https://www.grc.com/sqrl/protocol.htm
	// A SQRL client will only interpret a URL link beginning
	// with “sqrl://” or “qrl://” as a valid SQRL link.
)

// Cmd are the different commands a SQRL client
// can issue to the server.
type Cmd string

const (
	// CmdUnknown is used when a client's query was not understood.
	CmdUnknown = Cmd("")

	// CmdQuery allow the SQRL client to determine whether its user
	// is known to the SQRL server and to optionally obtain server-stored
	// data which the client may need.
	//
	// With each query command, the client assert its user's current identity
	// key and optionally a previous identity key.
	CmdQuery = Cmd("query")

	// CmdIdent requests will usually follow one or more query CmdQuery.
	//
	// Whereas the query queries allow the client to obtain information from
	// the server, the ident query requests the web server to accept the user's
	// identity assertion as it is provided by this signed query.
	CmdIdent = Cmd("ident")

	// CmdDisable instructs the web server to immediately disable the SQRL
	// system's authentication privilege for this domain.
	//
	// This might be requested if the user had reason to believe that their
	// current SQRL identity key had been compromised.
	CmdDisable = Cmd("disable")

	// CmdEnable is the reverse of the ‘disable’ query. It re-enables SQRL system
	// identity authentication for the user's account.
	//
	// Unlike ‘disable’, however, ‘enable’ requires the additional authorization
	// provided by the account's current unlock request signature (urs).
	CmdEnable = Cmd("enable")

	// CmdRemove instructs the web server to immediately remove all trace of this
	// SQRL identity from the server.
	//
	// For example, this process would allow an account to be disassociated from
	// one SQRL identity and subsequently reassociated with another.
	CmdRemove = Cmd("remove")
)

// Opt are used to indicate client preferences for a SQRL request.
type Opt string

const (
	OptNoIPTest = Opt("noiptest")
	OptSQRLOnly = Opt("sqrlonly")
	OptHardlock = Opt("hardlock")
	OptCPS      = Opt("cps")
	OptSUK      = Opt("suk")
)
