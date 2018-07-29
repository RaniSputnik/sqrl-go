package client

import (
	"encoding/base64"
	"strings"
)

func QueryCmd(idk string) string {
	vals := []string{
		"ver=1", // TODO: set version dynamically
		"cmd=query",
		"idk=" + idk,
		// TODO: "pidk=" + previousIdentityKey,
		// TODO: "suk=" + serverUnlockKey,
		// TODO: "vuk=" + verifyUnlockKey,

		"", // Must end with a final newline
	}
	return b64(strings.Join(vals, "\n"))
}

// b64 encodes some data to a string without padding
func b64(src string) string {
	encoder := base64.URLEncoding.WithPadding(base64.NoPadding)
	return encoder.EncodeToString([]byte(src))
}
