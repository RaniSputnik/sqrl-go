package sqrl

import (
	"errors"
	"fmt"
	"strings"
)

var errEmptyInput = errors.New("empty input")

func parseMsg(raw string) (map[string]string, error) {
	if raw == "" {
		return nil, errEmptyInput
	}

	bytes, err := Base64.DecodeString(raw)
	if err != nil {
		return nil, err
	}

	// TODO: avoid parsing ver parameter twice
	if !strings.HasPrefix(string(bytes), "ver=") {
		return nil, errors.New("must start with ver parameter")
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
		if _, ok := vals[key]; ok {
			return nil, fmt.Errorf("duplicate key '%s'", key)
		}
		vals[key] = val
	}

	return vals, nil
}

func parseVer(rawVer string) ([]string, error) {
	if rawVer == "" {
		return nil, errEmptyInput
	}
	return strings.Split(rawVer, ","), nil
}
