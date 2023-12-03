package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRand, err := rand.Read(b)

	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}

	if nRand < n {
		return nil, fmt.Errorf("bytes: didn't read enough bytes")
	}

	return b, nil
}

// String returns a random string using crypto/rand.
// n is the number of bytes being used to generate the string
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("string: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}