package parser

import (
	"artisons/conf"
	"errors"
	"fmt"
	"testing"
)

const image = "https://upload.wikimedia.org/wikipedia/commons/c/ca/1x1.png"

var line = []string{
	"123456", "Product", "product", "12.4", "EUR", "1", "online", "Best product", image, "164", "gifts;garden", "1234", "color:blue",
}

var header = []string{"sku", "title", "price", "currency", "quantity", "status", "description", "images", "weight", "tags", "links", "options"}

func TestImport(t *testing.T) {
	var tests = []struct {
		name   string
		header func(h []string) []string
		line   func(l []string) []string
		count  int
		err    error
	}{
		{
			"header missing",
			func(h []string) []string { return []string{} },
			func(l []string) []string { return l },
			0,
			errors.New("the csv is invalid"),
		},
		{
			"first cell invalid",
			func(h []string) []string { return []string{"id"} },
			func(l []string) []string { return l },
			0,
			errors.New("the csv is invalid"),
		},
		{
			"miss required fields",
			func(h []string) []string { return h },
			func(l []string) []string { return l[1:3] },
			0,
			nil,
		},
		{
			"bad sku",
			func(h []string) []string { return h },
			func(l []string) []string { l[0] = "sku!"; return l },
			0,
			nil,
		},
		{
			"title missing",
			func(h []string) []string { return h },
			func(l []string) []string { l[1] = ""; return l },
			0,
			nil,
		},
		{
			"slug missing",
			func(h []string) []string { return h },
			func(l []string) []string { l[2] = ""; return l },
			0,
			nil,
		},
		{
			"price missing",
			func(h []string) []string { return h },
			func(l []string) []string { l[3] = ""; return l },
			0,
			nil,
		},
		{
			"price invalid",
			func(h []string) []string { return h },
			func(l []string) []string { l[3] = "invalid"; return l },
			0,
			nil,
		},
		{
			"quantiy missing",
			func(h []string) []string { return h },
			func(l []string) []string { l[5] = ""; return l },
			0,
			nil,
		},
		{
			"quantiy invalid",
			func(h []string) []string { return h },
			func(l []string) []string { l[5] = "invalid"; return l },
			0,
			nil,
		},
		{
			"status invalid",
			func(h []string) []string { return h },
			func(l []string) []string { l[6] = "invalid"; return l },
			0,
			nil,
		},
		{
			"description missing",
			func(h []string) []string { return h },
			func(l []string) []string { l[7] = ""; return l },
			0,
			nil,
		},
		{
			"image missing",
			func(h []string) []string { return h },
			func(l []string) []string { l[8] = ""; return l },
			0,
			nil,
		},
		{
			"image invalid",
			func(h []string) []string { return h },
			func(l []string) []string { l[8] = "invalid"; return l },
			0,
			nil,
		},
		{
			"image not found",
			func(h []string) []string { return h },
			func(l []string) []string {
				l[8] = "https://upload.wikimedia.org/wikipedia/commons/c/ca/1x1_toto.png"
				return l
			},
			0,
			nil,
		},
		{
			"image bad extension",
			func(h []string) []string { return h },
			func(l []string) []string {
				l[8] = "https://upload.wikimedia.org/wikipedia/commons/c/ca/1x1.svg"
				return l
			},
			0,
			nil,
		},
		{
			"image local bad extension",
			func(h []string) []string { return h },
			func(l []string) []string {
				l[8] = "../../web/data/product.svg"
				return l
			},
			0,
			nil,
		},
		{
			"image local not found",
			func(h []string) []string { return h },
			func(l []string) []string {
				l[8] = "../../web/data/toto.png"
				return l
			},
			0,
			nil,
		},
		{
			"weight invalid",
			func(h []string) []string { return h },
			func(l []string) []string {
				l[9] = "toto"
				return l
			},
			0,
			nil,
		},
		{
			"options invalid",
			func(h []string) []string { return h },
			func(l []string) []string {
				l[12] = "toto"
				return l
			},
			0,
			nil,
		},
		{
			"count=1",
			func(h []string) []string { return h },
			func(l []string) []string {
				return l
			},
			1,
			nil,
		},
		{
			"local image",
			func(h []string) []string { return h },
			func(l []string) []string {
				l[8] = "../../web/data/product.png"
				return l
			},
			1,
			nil,
		},
		{
			"no option",
			func(h []string) []string { return h },
			func(l []string) []string {
				l[12] = ""
				return l
			},
			1,
			nil,
		},
		{
			"no link",
			func(h []string) []string { return h },
			func(l []string) []string {
				l[11] = ""
				return l
			},
			1,
			nil,
		},
		{
			"no weight",
			func(h []string) []string { return h },
			func(l []string) []string {
				l[9] = ""
				return l
			},
			1,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := make([]string, len(header))
			copy(h, header)

			l := make([]string, len(line))
			copy(l, line)
			csv := lines{tt.header(h), tt.line(l)}

			count, err := Import(csv, conf.DefaultMID)
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
				t.Fatalf(`err = %v, want %v`, err, tt.err)
			}

			if count != tt.count {
				t.Fatalf(`count = %d, want %d`, count, tt.count)
			}
		})
	}
}
