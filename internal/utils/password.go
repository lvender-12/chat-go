package utils

import (
	"crypto/rand"
	"crypto/subtle"

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

func HashingPassword(password string) ([]byte, error) {
	salt, err := generateSalt(saltLength)
	if err != nil {
		return nil, err
	}

	hashedPassword := argon2.IDKey([]byte(password), salt, timeParams, memory, threads, keyLength)
	return hashedPassword, nil
}

func VerifyPassword(password, hashPassword string) bool {
	newhashed, err := HashingPassword(password)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(newhashed, []byte(hashPassword)) == 1
}
