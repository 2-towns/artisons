// tests gather test utilites
package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func Context() context.Context {
	return context.WithValue(context.Background(), middleware.RequestIDKey, fmt.Sprintf("%d", time.Now().UnixMilli()))

}
