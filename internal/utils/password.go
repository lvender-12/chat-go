package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"

	"golang.org/x/crypto/argon2"
)

const (
	timeParams uint32 = 1
	memory     uint32 = 64 * 1024
	threads    uint8  = 4
	saltLength        = 16
	keyLength  uint32 = 32
)

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func HashingPassword(password string) (string, string, error) {
	salt, err := generateSalt(saltLength)
	if err != nil {
		return "", "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		timeParams,
		memory,
		threads,
		keyLength,
	)

	hashHex := hex.EncodeToString(hash)
	saltHex := hex.EncodeToString(salt)

	return hashHex, saltHex, nil
}

func VerifyPassword(password, hashHex, saltHex string) bool {
	hash, err := hex.DecodeString(hashHex)
	if err != nil {
		return false
	}

	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return false
	}

	newHash := argon2.IDKey(
		[]byte(password),
		salt,
		timeParams,
		memory,
		threads,
		keyLength,
	)

	return subtle.ConstantTimeCompare(newHash, hash) == 1
}
