// tests gather test utilites
package tests

import (
	"context"
	"fmt"
	"gifthub/http/contexts"
	"gifthub/string/stringutil"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/text/language"
)

func Context() context.Context {
	var ctx context.Context = context.WithValue(context.Background(), middleware.RequestIDKey, fmt.Sprintf("%d", time.Now().UnixMilli()))
	ctx = context.WithValue(ctx, contexts.Locale, language.English)

	rid, _ := stringutil.Random()
	ctx = context.WithValue(ctx, middleware.RequestIDKey, rid)

	return context.WithValue(ctx, contexts.Cart, fmt.Sprintf("%d", time.Now().UnixMilli()))
}
