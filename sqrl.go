package sqrl

import (
	"encoding/base64"
)

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
	// OptNoIPTest instructs the server to ignore any IP mismatch and to proceed
	// to process the client's query even if the IPs do not match. By default,
	// SQRL servers fail any incoming SQRL query whose IP does not match the IP
	// encoded into the query's nut.
	OptNoIPTest = Opt("noiptest")

	// OptSQRLOnly disables any alternative non-SQRL authentication capability.
	OptSQRLOnly = Opt("sqrlonly")

	// OptHardlock disables any alternative “out of band” change to this user's
	// SQRL identity, such as traditional and weak “what as your favorite pet's
	// name” non-SQRL identity authentication.
	OptHardlock = Opt("hardlock")

	// OptCPS informs the server that the client has established a secure and private
	// means of returning a server-supplied logged-in session URL to the web browser
	// after successful authentication.
	// https://www.grc.com/sqrl/semantics.htm
	OptCPS = Opt("cps")

	// OptSUK instructs the SQRL server to return the stored server unlock key
	// (SUK) associated with whichever identity matches the identity supplied
	// by the SQRL client.
	OptSUK = Opt("suk")
)

// TIF represent Transaction Interaction Flags sent from
// the server to provide information about a client request.
type TIF int

const (
	// TIFCurrentIDMatch indicates that the web server has found
	// an identity association for the user based upon the default
	// (current) identity credentials supplied by the client:
	// the IDentity Key (idk) and the IDentity Signature (ids).
	TIFCurrentIDMatch = TIF(0x01)

	// TIFPreviousIDMatch indicates that the web server has found
	// an identity association for the user based upon the previous
	// identity credentials supplied by the client in the previous
	// IDentity Key (pidk) and the previous IDentity Signature (pids).
	TIFPreviousIDMatch = TIF(0x02)

	// TIFIPMatch indicates that the IP address of the entity which
	// requested the initial logon web page containing the SQRL link URL
	// (and probably encoded into the SQRL link URL's “nut”) is the same
	// IP address from which the SQRL client's query was received for
	// this reply.
	TIFIPMatch = TIF(0x04)

	// TIFSQRLDisabled indicates that SQRL authentication for this
	// identity has previously been disabled. This bit can only be reset,
	// and the identity re-enabled for authentication, by the client issuing
	// an “enable” command signed by the unlock request signature (urs) for
	// the identity known to the server.
	TIFSQRLDisabled = TIF(0x08)

	// TIFFunctionNotSupported indicates that the client requested one or
	// more SQRL functions (through command verbs) that the server does
	// not currently support.
	TIFFunctionNotSupported = TIF(0x10)

	// TIFTransientError indicates that the client's signature(s) are correct,
	// but something about the client's query prevented the command from
	// completing. This is the server's way of instructing the client to retry
	// and reissue the immediately previous command using the fresh ‘nut=’
	// crypto material and ‘qry=’ url the server will have also just returned
	// in its reply.
	TIFTransientError = TIF(0x20)

	// TIFCommandFailed indicates that the web server has encountered a problem
	// fully processing the client's query. In any such case, no change will
	// be made to the user's account status.
	TIFCommandFailed = TIF(0x40)

	// TIFClientFailure is set by the server when some aspect of the client's
	// submitted query ‑ other than expired but otherwise valid transaction
	// state information ‑ was incorrect and prevented the server from
	// understanding and/or completing the requested action
	TIFClientFailure = TIF(0x80)

	// TIFBadIDAssociation is set by the server when a SQRL identity which
	// may be associated with the query nut does not match the SQRL ID used
	// to submit the query
	TIFBadIDAssociation = TIF(0x100)
)

// Base64 is the encoder that should be used when
// encoding keys, signatures and other payloads for
// transmission.
//
// It is the standard URL encoding with no padding to
// avoid confusion around the meaning of the = when
// form encoding SQRL requests.
var Base64 = base64.URLEncoding.WithPadding(base64.NoPadding)
