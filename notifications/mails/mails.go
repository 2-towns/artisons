// Package emails contains all the code related to sending emails
package mails

import (
	"gifthub/conf"
	"net/smtp"
)

// Send an email.
// Only text format is supported for now.
func Send(email, message string) error {
	from := conf.EmailUsername
	password := conf.EmailPassword

	to := []string{
		email,
	}

	smtpHost := conf.EmailHost
	smtpPort := conf.EmailPort

	m := []byte(message)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, m)
}
