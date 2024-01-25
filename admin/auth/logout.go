package auth

import (
	"gifthub/http/contexts"
	"gifthub/http/httperrors"
	"gifthub/users"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value(contexts.User).(users.User)

	if err := user.Logout(ctx); err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	w.Header().Set("HX-Redirect", "/auth/index.html")
	w.Write([]byte(""))
}
