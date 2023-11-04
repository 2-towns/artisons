// Package emails contains all the code related to sending emails
package mails

import (
	"context"
	"gifthub/conf"
	"log/slog"
	"net/smtp"
)

// Send an email.
// Only text format is supported for now.
func Send(c context.Context, email, message string) error {
	l := slog.With(slog.String("email", email))
	l.LogAttrs(c, slog.LevelInfo, "sending a new email")

	from := conf.EmailUsername
	password := conf.EmailPassword

	to := []string{
		email,
	}

	smtpHost := conf.EmailHost
	smtpPort := conf.EmailPort

	m := []byte(message)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, m)

	l.LogAttrs(c, slog.LevelInfo, "the email is sent")

	return err
}
