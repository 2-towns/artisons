package shops

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/validators"
	"log"
	"log/slog"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type Contact struct {
	Name    string `validate:"required"`
	Address string `validate:"required"`
	City    string `validate:"required"`
	Zipcode string `validate:"required"`
	Phone   string
	Email   string `validate:"required,email"`
	Logo    string `validate:"required"`
	Banner1 string
	Banner2 string
	Banner3 string
}

type Settings struct {
	Contact

	// Active defines if the store is available or not
	Active bool

	// Guest allows to accept guest order
	Guest bool

	// Show quantity in product page
	Quantity bool

	// Enable the stock managment
	Stock bool

	// Number of days during which the product is considered 'new'
	New bool

	// Max items per page
	Items int

	// Mininimum order
	Min int

	// Redirect after  the product was added to the cart
	Redirect bool

	// Display last products when the quantity is under the amount.
	// Set to zero to disable this feature.
	LastProducts int

	// AdvancedSearch enables the advanced search
	AdvancedSearch bool

	// Cache enables the advanced search
	Cache bool

	// Google map key used for geolocation api
	GmapKey string
}

var Data Settings

func init() {
	ctx := context.Background()
	d, err := db.Redis.HGetAll(ctx, "shop").Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get shop info", slog.String("error", err.Error()))
		log.Fatalln(err)
	}

	var items int = 0
	if d["items"] != "" {
		i, err := strconv.ParseInt(d["items"], 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the items", slog.String("items", d["items"]), slog.String("error", err.Error()))
		} else {
			items = int(i)
		}
	}

	var min int = 0
	if d["min"] != "" {
		i, err := strconv.ParseInt(d["min"], 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the min", slog.String("min", d["min"]), slog.String("error", err.Error()))
		} else {
			min = int(i)
		}
	}

	var las int = 0
	if d["milast_productsn"] != "" {
		i, err := strconv.ParseInt(d["last_products"], 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the last_products", slog.String("last_products", d["last_products"]), slog.String("error", err.Error()))
		} else {
			las = int(i)
		}
	}

	Data = Settings{
		Contact: Contact{
			Name:    d["name"],
			Address: d["address"],
			City:    d["city"],
			Zipcode: d["zipcode"],
			Phone:   d["phone"],
			Email:   d["email"],
			Logo:    d["logo"],
			Banner1: d["banner_1"],
			Banner2: d["banner_2"],
			Banner3: d["banner_3"],
		},
		Guest:          d["guest"] == "1",
		Quantity:       d["quantity"] == "1",
		Stock:          d["stock"] == "1",
		New:            d["new"] == "1",
		Items:          items,
		Min:            min,
		Redirect:       d["redirect"] == "1",
		LastProducts:   las,
		AdvancedSearch: d["advanced_search"] == "1",
		Cache:          d["cache"] == "1",
		GmapKey:        d["gmap_key"],
	}
}

func (s Contact) Validate(c context.Context) error {
	slog.LogAttrs(c, slog.LevelInfo, "validating a contact settings")

	if err := validators.V.Struct(s); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot validate the shop", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input_%s_invalid", low)
	}

	return nil
}

func (s Settings) Save(c context.Context) error {
	l := slog.With(slog.String("name", s.Name))
	l.LogAttrs(c, slog.LevelInfo, "trying to save the shop")

	guest := "0"
	if s.Guest {
		guest = "1"
	}

	quantity := "0"
	if s.Quantity {
		quantity = "1"
	}

	stock := "0"
	if s.Stock {
		stock = "1"
	}

	new := "0"
	if s.New {
		new = "1"
	}

	redirect := "0"
	if s.Redirect {
		redirect = "1"
	}

	asearch := "0"
	if s.AdvancedSearch {
		asearch = "1"
	}

	cache := "0"
	if s.Cache {
		cache = "1"
	}

	now := time.Now()
	_, err := db.Redis.HSet(context.Background(), "shop",
		"logo", path.Join(conf.ImgProxy.Path, s.Logo),
		"name", s.Name,
		"address", s.Address,
		"city", s.City,
		"zipcode", s.Zipcode,
		"phone", s.Phone,
		"email", s.Email,
		"logo", s.Logo,
		"guest", guest,
		"quantity", quantity,
		"stock", stock,
		"new", new,
		"items", fmt.Sprintf("%d", s.Items),
		"min", fmt.Sprintf("%d", s.Min),
		"redirect", redirect,
		"last_products", fmt.Sprintf("%d", s.LastProducts),
		"advanced_search", asearch,
		"cache", cache,
		"gmap_key", s.GmapKey,
		"banner_1", s.Banner1,
		"banner_2", s.Banner2,
		"banner_3", s.Banner3,
		"updated_at", now.Unix(),
	).Result()

	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot save the shop", slog.String("error", err.Error()))
		return errors.New("error_http_general")
	}

	return nil
}
