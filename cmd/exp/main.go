package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	secret := "MySuperSecretPhrase"
	password := "ThisIsMyPassword"

	h := hmac.New(sha256.New, []byte(secret))

	h.Write([]byte(password))

	output := h.Sum(nil)

	fmt.Println(hex.EncodeToString(output))
	// => 827e2efb277ea22df6e9559ecd5dd5448b7da6f1ba3d63fde0f14b91714e9bb7
}
