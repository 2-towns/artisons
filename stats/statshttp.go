package stats

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/cookies"
	"artisons/http/httperrors"
	"artisons/string/stringutil"
	"artisons/tracking"
	"artisons/users"
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mileusna/useragent"
)

func Demo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(contexts.User).(users.User)
	if !ok {
		httperrors.Catch(w, ctx, "something_went_wrong", 400)
		return
	}

	_, err := user.ToggleDemo(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	w.Header().Set("HX-Redirect", "/admin/index.html")
	w.Write([]byte(""))
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var did string
		c, err := r.Cookie(cookies.Device)
		if err != nil || c.Value == "" {
			did, err = stringutil.Random()
			if err != nil {
				slog.LogAttrs(r.Context(), slog.LevelError, "cannot generate device id", slog.String("error", err.Error()))
				httperrors.Page(w, r.Context(), err.Error(), 500)
				return
			}
		} else {
			did = c.Value
		}

		cookie := &http.Cookie{
			Name:     cookies.Device,
			Value:    did,
			MaxAge:   int(conf.Cookie.MaxAge),
			Path:     "/",
			HttpOnly: true,
			Secure:   conf.Cookie.Secure,
			Domain:   conf.Cookie.Domain,
		}

		http.SetCookie(w, cookie)

		ctx := context.WithValue(r.Context(), contexts.Device, did)

		ua := useragent.Parse(r.Header.Get("User-Agent"))
		go Visit(ctx, ua, VisitData{
			URL:     r.URL.Path,
			Referer: r.Referer(),
		})

		if conf.EnableTrackingLog {
			data := map[string]string{
				"url":     r.URL.Path,
				"referer": fmt.Sprintf("'%s'", r.Referer()),
				"ua":      fmt.Sprintf("'%s'", r.Header.Get("User-agent")),
			}

			go tracking.Log(ctx, "access", data)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
