package security

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateTokenHash() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
