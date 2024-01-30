package users

import (
	"artisons/db"
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Session struct {
	ID      string
	Device  string
	WPToken string
	TTL     time.Duration
}

// Sessions retrieve the active user sessions.
func (u User) Sessions(ctx context.Context) ([]Session, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "searching sessions", slog.Int("uid", u.ID))

	qs := fmt.Sprintf("FT.SEARCH %s @type:{session}@uid:{%d}", db.SessionIdx, u.ID)
	qs += " SORTBY updated_at desc LIMIT 0 9999 DIALECT 2"

	slog.LogAttrs(ctx, slog.LevelInfo, "preparing redis request", slog.String("query", qs))

	args, err := db.SplitQuery(ctx, qs)
	if err != nil {
		return []Session{}, err
	}

	cmds, err := db.Redis.Do(ctx, args...).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot run the search query", slog.String("error", err.Error()))
		return []Session{}, err
	}

	res := cmds.(map[interface{}]interface{})
	total := res["total_results"].(int64)
	results := res["results"].([]interface{})
	sessions := []Session{}

	for _, value := range results {
		m := value.(map[interface{}]interface{})
		attributes := m["extra_attributes"].(map[interface{}]interface{})
		data := db.ConvertMap(attributes)

		ttl, err := db.Redis.TTL(ctx, "session:"+data["id"]).Result()
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the session ttl", slog.String("id", data["id"]), slog.String("error", err.Error()))
			continue
		}

		session := Session{
			ID:      data["id"],
			Device:  data["device"],
			WPToken: data["wptoken"],
			TTL:     ttl,
		}

		sessions = append(sessions, session)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "search done", slog.Int64("results", total))

	return sessions, nil
}
