package mails

import (
	"gifthub/tests"
	"testing"

	"github.com/go-faker/faker/v4"
)

func TestSendReturnsNilWhenSuccess(t *testing.T) {
	email := faker.Email()
	c := tests.Context()
	if err := Send(c, email, "subject", faker.Sentence()); err != nil {
		t.Fatalf("Send(c, email, faker.Sentence()) = %v, want nil", err)
	}
}
