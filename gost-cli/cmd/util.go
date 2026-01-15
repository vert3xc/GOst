package cmd

import (
	"encoding/hex"
	"errors"
)

func decodeHexArg(s string) ([]byte, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, errors.New("invalid hex input")
	}
	return b, nil
}