// Package emails contains all the code related to sending emails
package mails

import "net/smtp"

func Send(email string, message string) error {
	// Sender data.
	from := "a3a5f2d396a820"
	password := "12fcfd3c6edb95"

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "sandbox.smtp.mailtrap.io"
	smtpPort := "587"

	// Message.
	m := []byte(message)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, m)
}
