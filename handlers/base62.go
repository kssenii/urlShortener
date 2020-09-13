package handlers

import (
	"errors"
	"fmt"
	"strings"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length   = int64(len(alphabet))
)

func EncodeBase62(data *Data) {

	if data.ID == 0 {
		data.ShortURL = string(alphabet[0])
	}

	for n := data.ID; n > 0; n = n / length {
		data.ShortURL = string(alphabet[n%length]) + data.ShortURL
	}
}

// Decode converts a base62 token to int.
func DecodeBase62(data *Data) error {

	var n int64
	for _, c := range []byte(data.ShortURL) {

		pos := strings.IndexByte(alphabet, c)
		if pos < 0 {
			return errors.New(fmt.Sprintf("Unexpected character %c at position %d", c, pos))
		}

		n = length*n + int64(pos)
	}

	data.ID = n
	return nil
}
