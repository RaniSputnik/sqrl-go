package sqrl

import (
	"fmt"
	"strings"
)

func parseMsg(raw string) (map[string]string, error) {
	bytes, err := Base64.DecodeString(raw)
	if err != nil {
		return nil, err
	}

	form := strings.Split(string(bytes), "\n")
	vals := map[string]string{}
	for _, keyval := range form {
		if keyval == "" {
			continue
		}

		pair := strings.SplitN(keyval, "=", 2)
		if len(pair) < 2 {
			return nil, fmt.Errorf("invalid value '%s', should be in the form: key=value\\n", keyval)
		}
		key := strings.TrimSpace(pair[0])
		val := strings.TrimSpace(pair[1])
		vals[key] = val
	}

	return vals, nil
}

func parseVer(rawVer string) []string {
	return strings.Split(rawVer, ",")
}
