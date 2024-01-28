package users

import (
	"context"
	"errors"
	"fmt"
	"artisons/db"
	"log/slog"
)

// AddWPToken registers a vapid webpush token
// to receive push notifications.
// The data are stored with:
// - auth:sid:session wptoken => the wptoken related to the session
// The token is the string representation of the
// JSON token.
func (u User) AddWPToken(ctx context.Context, token string) error {
	l := slog.With(slog.String("token", token))
	l.LogAttrs(ctx, slog.LevelInfo, "adding a new webpush token")

	if token == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the token")
		return errors.New("input:wptoken")
	}

	if _, err := db.Redis.HSet(ctx, fmt.Sprintf("auth:%s:session", u.SID), "wptoken", token).Result(); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot store the token", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "token stored successfully")

	return nil
}

// DeleteWPToken removes a vapid webpush token linked to
// a session.
// The data are stored with:
// - auth:sid:session wptoken => the wptoken related to the session
func (u User) DeleteWPToken(ctx context.Context, sid string) error {
	l := slog.With(slog.String("sid", sid))
	l.LogAttrs(ctx, slog.LevelInfo, "deleting a web push token")

	if sid == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the session id")
		return errors.New("your are not authorized to process this request")
	}

	if _, err := db.Redis.HDel(ctx, fmt.Sprintf("auth:%s:session", u.SID), "wptoken").Result(); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot delete the token from redis", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "token deleted successfully")

	return nil
}
