package main

import "crypto/rand"

const (
	alphaChars    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	numChars      = "0123456789"
	alphaNumChars = alphaChars + numChars
)

// generateToken generates a cryptographically random,
// alphanumeric string of length n.
func generateToken(totalLen int) (string, error) {
	bytes := make([]byte, totalLen)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for k, v := range bytes {
		bytes[k] = alphaNumChars[v%byte(len(alphaNumChars))]
	}
	return string(bytes), nil
}
