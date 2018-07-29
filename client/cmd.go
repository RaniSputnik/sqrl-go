package client

import (
	sqrl "github.com/RaniSputnik/sqrl-go"
)

func QueryCmd(idk sqrl.Identity) string {
	query := sqrl.ClientMsg{
		Ver: v1Only,
		Cmd: sqrl.CmdQuery,
		Idk: idk,
		// TODO: "pidk=" + previousIdentityKey,
		// TODO: "suk=" + serverUnlockKey,
		// TODO: "vuk=" + verifyUnlockKey,
	}

	// Swallow error, only occurs if we have
	// an incomplete client msg
	val, _ := query.Encode()
	return val
}

var v1Only = []string{sqrl.V1}
