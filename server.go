package sqrl

// ServerMsg is either the base64 encoded sqrl:// URL
// sent by the server to initiate the exchange OR
// the exact value of the 'server' parameter from the
// previous transaction.
type ServerMsg string
