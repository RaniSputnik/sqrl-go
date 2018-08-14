package sqrl

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
