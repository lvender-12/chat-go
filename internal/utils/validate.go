package utils

import (
	"net/mail"
	"strings"
)

func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)

	if email == "" {
		return false
	}

	address, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	return address.Address == email
}
