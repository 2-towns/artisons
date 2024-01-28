package addresses

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"artisons/conf"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

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
