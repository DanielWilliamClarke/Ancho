package deveui

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateHexCode generates a random hex value of length.
func GenerateHexCode(length int) (string, error) {
	// to reduce the amount of memory allocated divide the length by 2
	// as encoded hex string are of length * 2
	bytes := make([]byte, (length+1)/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}
