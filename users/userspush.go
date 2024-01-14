package users

import (
	"context"
	"errors"
	"fmt"
	"gifthub/db"
	"log/slog"
)

// AddWPToken registers a vapid webpush token
// to receive push notifications.
// The data are stored with:
// - auth:sid:session wptoken => the wptoken related to the session
// The token is the string representation of the
// JSON token.
func (u User) AddWPToken(c context.Context, token string) error {
	l := slog.With(slog.String("token", token))
	l.LogAttrs(c, slog.LevelInfo, "adding a new webpush token")

	if token == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the token")
		return errors.New("input:wptoken")
	}

	ctx := context.Background()
	if _, err := db.Redis.HSet(ctx, fmt.Sprintf("auth:%s:session", u.SID), "wptoken", token).Result(); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot store the token", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "token stored successfully")

	return nil
}

// DeleteWPToken removes a vapid webpush token linked to
// a session.
// The data are stored with:
// - auth:sid:session wptoken => the wptoken related to the session
func (u User) DeleteWPToken(c context.Context, sid string) error {
	l := slog.With(slog.String("sid", sid))
	l.LogAttrs(c, slog.LevelInfo, "deleting a web push token")

	if sid == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the session id")
		return errors.New("your are not authorized to process this request")
	}

	ctx := context.Background()
	if _, err := db.Redis.HDel(ctx, fmt.Sprintf("auth:%s:session", u.SID), "wptoken").Result(); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot delete the token from redis", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "token deleted successfully")

	return nil
}
