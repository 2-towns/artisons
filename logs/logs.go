// logs provide the utilities for logging
package logs

import (
	"context"
	"artisons/http/contexts"
	"artisons/users"
	"log"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5/middleware"
)

type RequestIDHandler struct {
	slog.Handler
}

func (h RequestIDHandler) Handle(ctx context.Context, r slog.Record) error {
	if rid, ok := ctx.Value(middleware.RequestIDKey).(string); ok {
		r.Add("request_id", slog.StringValue(rid))
	}

	if u, ok := ctx.Value(contexts.User).(users.User); ok {
		r.Add("user_id", slog.IntValue(u.ID))
	}

	return h.Handler.Handle(ctx, r)
}

func Init() {
	//handler := RequestIDHandler{slog.Default().Handler()}
	handler := RequestIDHandler{slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
	})}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	// https://github.com/golang/go/issues/61892#issuecomment-1675123776
	log.SetOutput(os.Stderr)
}
