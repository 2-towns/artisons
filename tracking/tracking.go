// tracking is reponsible to keep trace of the application in logs
// in order to improve the statistics data in the future
package tracking

import (
	"artisons/conf"
	"artisons/http/contexts"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/text/language"
)

func Log(ctx context.Context, action string, data map[string]string) error {
	tracking, ok := ctx.Value(contexts.Tracking).(bool)
	if !tracking || !ok {
		return nil
	}

	l := slog.With(slog.String("action", action))
	l.LogAttrs(ctx, slog.LevelInfo, "writing tracking log")

	folder := conf.WorkingSpace + "web/tracking"
	now := time.Now()
	name := fmt.Sprintf("tracking-%s.log", now.Format("20060102"))
	f, err := os.OpenFile(path.Join(folder, name),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "error when opening tracking log file", slog.String("error", err.Error()))
		return err
	}

	rid := ctx.Value(contexts.RequestID).(string)
	cid := ctx.Value(contexts.Device).(string)
	lang := ctx.Value(contexts.Locale).(language.Tag)
	parts := []string{
		fmt.Sprintf("time:%d", now.Unix()),
		fmt.Sprintf("rid:%s", rid),
		fmt.Sprintf("cid:%s", cid),
		fmt.Sprintf("lang:%s", lang),
	}

	uid, ok := ctx.Value(contexts.UserID).(int)
	if ok && uid > 0 {
		parts = append(parts, fmt.Sprintf("uid:%d", uid))
	}

	for key, value := range data {
		parts = append(parts, fmt.Sprintf("%s:%s", key, value))
	}

	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s\n", strings.Join(parts, " "))); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "error when writing tracking log file", slog.String("error", err.Error()))
		return err
	}

	return nil
}
