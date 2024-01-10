package admin

import (
	"fmt"
	"gifthub/blogs"
	"gifthub/http/httperrors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func DeleteBlog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "error_http_general")
		return
	}

	id := chi.URLParam(r, "id")

	iid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.String("id", id), slog.String("error", err.Error()))
		httperrors.Page(w, ctx, "error_http_blognotfound", 404)
	}

	err = blogs.Delete(ctx, iid)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	page := r.FormValue("page")
	query := r.FormValue("query")
	r.URL.RawQuery = fmt.Sprintf("page%s&query=%s", page, query)

	Blog(w, r)
}
