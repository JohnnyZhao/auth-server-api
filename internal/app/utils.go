package app

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"unicode"
)

func isAlphanumeric(s string) bool {
	alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return alphanumeric.MatchString(s)
}

func isASCIIString(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 || unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func isValidChar(s string) bool {
	for _, char := range s {
		if unicode.IsControl(char) {
			return false
		}
	}
	return true
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
