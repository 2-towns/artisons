package stats

import (
	"fmt"
	"gifthub/admin/urls"
	"gifthub/http/contexts"
	"gifthub/http/httperrors"
	"gifthub/tracking"
	"gifthub/users"
	"net/http"
	"strings"

	"github.com/mileusna/useragent"
)

func Demo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(contexts.User).(users.User)
	if !ok {
		httperrors.Catch(w, ctx, "something_went_wrong")
		return
	}

	_, err := user.ToggleDemo(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error())
		return
	}

	w.Header().Set("HX-Redirect", urls.AdminPrefix)
	w.Write([]byte(""))
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/public") || strings.HasPrefix(r.URL.Path, urls.AdminPrefix) || strings.HasPrefix(r.URL.Path, urls.AuthPrefix) {
			next.ServeHTTP(w, r.WithContext(r.Context()))
			return
		}

		ua := useragent.Parse(r.Header.Get("User-Agent"))

		go Visit(r.Context(), ua, VisitData{
			URL:     r.URL.Path,
			Referer: r.Referer(),
		})

		data := map[string]string{
			"url":     r.URL.Path,
			"referer": fmt.Sprintf("'%s'", r.Referer()),
			"ua":      fmt.Sprintf("'%s'", r.Header.Get("User-agent")),
		}

		go tracking.Log(r.Context(), "access", data)

		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}
