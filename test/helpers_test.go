package test

import (
	"testing"

	"chat-go/internal/utils"
)

func mustToken(t *testing.T, userID uint64, secret []byte) (string, error) {
	t.Helper()
	return utils.GenerateToken(userID, secret, 1)
}
