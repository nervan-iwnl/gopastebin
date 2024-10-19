package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomPasteLink() string {
	const idLength = 8

	randomBytes := make([]byte, idLength)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}

	pasteLink := hex.EncodeToString(randomBytes)

	return pasteLink
}
