package products

import (
	"artisons/db"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func Wish(ctx context.Context, uid int, pid string) error {
	l := slog.With(slog.String("pid", pid))
	l.LogAttrs(ctx, slog.LevelInfo, "adding to wish list")

	if _, err := db.Redis.ZAdd(ctx, fmt.Sprintf("wish:%d", uid), redis.Z{
		Member: pid,
		Score:  float64(time.Now().Unix()),
	}).Result(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot add to the wish list", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	return nil
}

func UnWish(ctx context.Context, uid int, pid string) error {
	l := slog.With(slog.String("pid", pid))
	l.LogAttrs(ctx, slog.LevelInfo, "removing from wish list")

	if _, err := db.Redis.ZRem(ctx, fmt.Sprintf("wish:%d", uid), pid).Result(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot remove from the wish list", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	return nil
}

func Wishes(ctx context.Context, uid int) ([]string, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "listing wish list")

	wishes, err := db.Redis.ZRange(ctx, fmt.Sprintf("wish:%d", uid), 0, 9999).Result()

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the wish list", slog.String("error", err.Error()))
		return []string{}, errors.New("something went wrong")
	}

	return wishes, nil
}

func HasWish(ctx context.Context, uid int, pid string) bool {
	slog.LogAttrs(ctx, slog.LevelInfo, "listing wish list")

	score, err := db.Redis.ZScore(ctx, fmt.Sprintf("wish:%d", uid), pid).Result()

	return err == nil && score > 0
}
