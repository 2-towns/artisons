// tests gather test utilites
package tests

import (
	"artisons/http/contexts"
	"artisons/string/stringutil"
	"context"
	"fmt"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/text/language"
)

func Context() context.Context {
	var ctx context.Context = context.WithValue(context.Background(), middleware.RequestIDKey, fmt.Sprintf("%d", time.Now().UnixMilli()))
	ctx = context.WithValue(ctx, contexts.Locale, language.English)

	rid, _ := stringutil.Random()
	ctx = context.WithValue(ctx, middleware.RequestIDKey, rid)
	ctx = context.WithValue(ctx, contexts.Demo, true)

	return context.WithValue(ctx, contexts.Device, fmt.Sprintf("%d", time.Now().UnixMilli()))
}
