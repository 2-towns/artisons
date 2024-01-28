package users

import (
	"context"
	"errors"
	"fmt"
	"artisons/db"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func (u User) Wish(ctx context.Context, pid string) error {
	l := slog.With(slog.String("pid", pid))
	l.LogAttrs(ctx, slog.LevelInfo, "adding to wish list")

	if _, err := db.Redis.ZAdd(ctx, fmt.Sprintf("wish:%d", u.ID), redis.Z{
		Member: pid,
		Score:  float64(time.Now().Unix()),
	}).Result(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot add to the wish list", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	return nil
}

func (u User) UnWish(ctx context.Context, pid string) error {
	l := slog.With(slog.String("pid", pid))
	l.LogAttrs(ctx, slog.LevelInfo, "removing from wish list")

	if _, err := db.Redis.ZRem(ctx, fmt.Sprintf("wish:%d", u.ID), pid).Result(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot remove from the wish list", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	return nil
}

func (u User) Wishes(ctx context.Context) ([]string, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "listing wish list")

	wishes, err := db.Redis.ZRange(ctx, fmt.Sprintf("wish:%d", u.ID), 0, 9999).Result()

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the wish list", slog.String("error", err.Error()))
		return []string{}, errors.New("something went wrong")
	}

	return wishes, nil
}

func (u User) HasWish(ctx context.Context, pid string) bool {
	slog.LogAttrs(ctx, slog.LevelInfo, "listing wish list")

	score, err := db.Redis.ZScore(ctx, fmt.Sprintf("wish:%d", u.ID), pid).Result()

	return err == nil && score > 0
}
