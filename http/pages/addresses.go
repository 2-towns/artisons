package pages

import (
	"gifthub/addresses"
	"gifthub/conf"
	"gifthub/http/httperrors"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
)

var addressesTpl *template.Template

func init() {
	var err error

	addressesTpl, err = templates.Build("addresses.html").ParseFiles(
		conf.WorkingSpace + "web/views/addresses.html",
	)

	if err != nil {
		log.Panicln(err)
	}
}

func Addresses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query().Get("q")

	addresses, err := addresses.Get(ctx, q, 10)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	data := struct {
		Data []string
	}{
		addresses,
	}

	w.Header().Set("Content-Type", "text/html")

	if err = addressesTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
