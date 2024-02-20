package auth

import (
	"artisons/http/cookies"
	"artisons/http/httperrors"
	"artisons/http/httphelpers"
	"artisons/users"
	"log/slog"
	"net/http"
	"strings"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sid, err := r.Cookie(cookies.SessionID)
	if err == nil {
		user := users.User{SID: sid.Value}
		err := user.Logout(ctx)

		cookie := httphelpers.NewCookie(cookies.SessionID, user.SID, -1)
		http.SetCookie(w, &cookie)

		if err != nil {
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}
	} else {
		slog.LogAttrs(ctx, slog.LevelInfo, "the user is not in context")
	}

	if strings.HasPrefix(r.Header.Get("HX-Current-Url"), "/admin") {
		w.Header().Set("HX-Redirect", "/sso")
		w.Write([]byte(""))
	} else {
		w.Header().Set("HX-Redirect", "/")
		w.Write([]byte(""))
	}
}
