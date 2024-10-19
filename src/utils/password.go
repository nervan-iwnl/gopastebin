package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func GenerateSalt() string {
	salt := make([]byte, 16)
	rand.Read(salt)
	return base64.StdEncoding.EncodeToString(salt)
}

func HashPassword(password string, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return fmt.Sprintf("%s$%s", salt, base64.StdEncoding.EncodeToString(hash[:]))
}

func CheckPassword(password string, hashedPassword string) bool {
	salt := strings.Split(hashedPassword, "$")[0]
	return hashedPassword == HashPassword(password, salt)
}
