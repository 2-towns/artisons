package addresses

import (
	"artisons/conf"
	"artisons/http/httperrors"
	"artisons/templates"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
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

type geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type properties struct {
	Label       string  `json:"label"`
	Score       float64 `json:"score"`
	HouseNumber string  `json:"housenumber"`
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	PostCode    string  `json:"postcode"`
	CityCode    string  `json:"citycode"`
	C           float64 `json:"x"`
	Y           float64 `json:"y"`
	City        string  `json:"city"`
	Context     string  `json:"context"`
	Importance  string  `json:"importance"`
	Street      string  `json:"street"`
}

type feature struct {
	Type        string     `json:"type"`
	Geometry    geometry   `json:"geometry"`
	Properties  properties `json:"properties"`
	Attribution string     `json:"attribution"`
	Licence     string     `json:"licence"`
	Query       string     `json:"query"`
	Limit       int        `json:"limit"`
}

type response struct {
	Type     string    `json:"type"`
	Version  string    `json:"version"`
	Features []feature `json:"features"`
}

func Get(ctx context.Context, pattern string, limit int) ([]string, error) {
	l := slog.With(slog.String("pattern", pattern), slog.Int("limit", limit))
	res, err := http.Get(fmt.Sprintf("%s/search/?q=%s&limit=%d", conf.AddressesFrApi, url.QueryEscape(pattern), limit))

	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot search for addresses")
		return []string{}, errors.New("something went wrong")
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot read the addresses response")
		return []string{}, errors.New("something went wrong")
	}

	var r response
	json.Unmarshal(data, &r)

	addresses := []string{}
	for _, val := range r.Features {
		addresses = append(addresses, val.Properties.Label)
	}

	return addresses, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query().Get("q")

	var addresses []string = []string{}

	if len(q) < 3 {
		slog.LogAttrs(ctx, slog.LevelInfo, "the query is too short", slog.String("q", q))
	} else {
		var err error
		addresses, err = Get(ctx, q, 10)
		if err != nil {
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}
	}

	data := struct {
		Data []string
	}{
		addresses,
	}

	w.Header().Set("Content-Type", "text/html")

	if err := addressesTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
