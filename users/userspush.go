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
// The token is the string representation of the
// JSON token.
func (u User) AddWPToken(c context.Context, token string) error {
	l := slog.With(slog.String("token", token))
	l.LogAttrs(c, slog.LevelInfo, "adding a new webpush token")

	if token == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the token")
		return errors.New("user_wptoken_required")
	}

	ctx := context.Background()
	if _, err := db.Redis.HSet(ctx, fmt.Sprintf("user:%d", u.ID), "wptoken:"+u.SID, token).Result(); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot store the token", slog.String("error", err.Error()))
		return errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "token stored successfully")

	return nil
}

// DeleteWPToken removes a vapid webpush token linked to
// a session.
func (u User) DeleteWPToken(c context.Context, sid string) error {
	l := slog.With(slog.String("sid", sid))
	l.LogAttrs(c, slog.LevelInfo, "deleting a web push token")

	if sid == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the session id")
		return errors.New("unauthorized")
	}

	ctx := context.Background()
	if _, err := db.Redis.HDel(ctx, fmt.Sprintf("user:%d", u.ID), "wptoken:"+u.SID).Result(); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot delete the token from redis", slog.String("error", err.Error()))
		return errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "token deleted successfully")

	return nil
}
