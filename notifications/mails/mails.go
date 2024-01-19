// Package emails contains all the code related to sending emails
package mails

import (
	"context"
	"fmt"
	"gifthub/conf"
	"gifthub/string/stringutil"
	"log/slog"
	"net/smtp"
	"time"
)

const rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700"

// Send an email.
// Only text format is supported for now.
func Send(ctx context.Context, email, subject, message string) error {
	l := slog.With(slog.String("email", email))
	l.LogAttrs(ctx, slog.LevelInfo, "sending a new email")

	mid, err := stringutil.Random()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot generated the message-id", slog.String("error", err.Error()))
		return err
	}

	from := conf.Email.From
	host := conf.Email.Host
	port := conf.Email.Port
	login := conf.Email.Username
	pass := conf.Email.Password
	to := []string{
		email,
	}
	msg := fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", to[0]) +
		fmt.Sprintf("Date: %s\r\n", time.Now().Format(rfc2822)) +
		fmt.Sprintf("Content-Type: %s\r\n", "text/plain; charset=us-ascii") +
		fmt.Sprintf("Message-ID: <%s@%s>\r\n", mid, conf.Email.Domain) +
		fmt.Sprintf("MIME-Version: %s\r\n", "1.0") +
		fmt.Sprintf("Subject: %s\r\n\r\n", subject) +
		fmt.Sprintf("%s\r\n", message)

	if conf.Email.Dry {
		slog.LogAttrs(ctx, slog.LevelInfo, "will not send email because of dry")
		return nil
	}

	auth := smtp.PlainAuth("", login, pass, host)
	err = smtp.SendMail(fmt.Sprintf("%s:%s", host, port), auth, from, to, []byte(msg))

	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot send the email", slog.String("error", err.Error()))
	} else {
		l.LogAttrs(ctx, slog.LevelInfo, "the email is sent")
	}

	return err
}
