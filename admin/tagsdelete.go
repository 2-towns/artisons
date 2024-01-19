package admin

import (
	"fmt"
	"gifthub/http/httperrors"
	"gifthub/tags"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	pid := chi.URLParam(r, "id")

	err := tags.Delete(ctx, pid)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	page := r.FormValue("page")
	r.URL.RawQuery = fmt.Sprintf("page%s", page)

	Tags(w, r)
}
