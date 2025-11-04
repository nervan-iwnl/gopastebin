package slug

import (
	"crypto/rand"
	"encoding/hex"
)

func Generate() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
