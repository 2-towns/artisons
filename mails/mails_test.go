package mails

import (
	"testing"

	"github.com/go-faker/faker/v4"
)

// TestSend send a test email
func TestSend(t *testing.T) {
	email := faker.Email()
	if err := Send(email, faker.Sentence()); err != nil {
		t.Fatalf("Send(email,faker.Sentence()), %v, nil, error", err)
	}
}
