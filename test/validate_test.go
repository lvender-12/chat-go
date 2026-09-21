package test

import (
	"testing"

	"chat-go/internal/utils"
)

func TestIsValidEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{name: "valid", email: "user@example.com", want: true},
		{name: "valid with plus", email: "user+tag@example.com", want: true},
		{name: "empty", email: "", want: false},
		{name: "missing at", email: "userexample.com", want: false},
		{name: "missing domain", email: "user@", want: false},
		{name: "spaces", email: "user @example.com", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := utils.IsValidEmail(tt.email)
			if got != tt.want {
				t.Fatalf("IsValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}
