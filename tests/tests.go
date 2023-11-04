// tests gather test utilites
package tests

import (
	"context"
	"fmt"
	"gifthub/locales"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/text/language"
)

func Context() context.Context {
	var ctx context.Context = context.WithValue(context.Background(), middleware.RequestIDKey, fmt.Sprintf("%d", time.Now().UnixMilli()))
	return context.WithValue(ctx, locales.ContextKey, language.English)

}
