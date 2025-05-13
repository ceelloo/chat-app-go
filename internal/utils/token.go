package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("Failed to generate token")
	}

	return base64.URLEncoding.EncodeToString(bytes)
}