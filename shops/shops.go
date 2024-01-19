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
	Logo    string
	Banner1 string
	Banner2 string
	Banner3 string
}

type ShopSettings struct {
	// Active defines if the store is available or not
	Active bool

	// Guest allows to accept guest order
	Guest bool

	// Show quantity in product page
	Quantity bool

	// Number of days during which the product is considered 'new'
	New bool

	// Max items per page
	Items int

	// Mininimum order
	Min int

	// Redirect after  the product was added to the cart
	Redirect bool

	// Cache enables the advanced search
	Cache bool

	// Google map key used for geolocation api
	GmapKey string

	// Enable the Redis fuzzy search
	FuzzySearch bool

	// If true, the search will look for the exact match pattern by default
	ExactMatchSearch bool

	// The brand color
	Color string

	// The default image width
	ImageWidth int

	// The default image height
	ImageHeight int
}

type Settings struct {
	Contact
	ShopSettings
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

	var width int = 0
	if d["image_width"] != "" {
		i, err := strconv.ParseInt(d["image_width"], 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the image_width", slog.String("image_width", d["image_width"]), slog.String("error", err.Error()))
		} else {
			width = int(i)
		}
	}

	var height int = 0
	if d["image_width"] != "" {
		i, err := strconv.ParseInt(d["image_height"], 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the image_height", slog.String("image_height", d["image_width"]), slog.String("error", err.Error()))
		} else {
			height = int(i)
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
		ShopSettings: ShopSettings{
			Guest:            d["guest"] == "1",
			Quantity:         d["quantity"] == "1",
			New:              d["new"] == "1",
			Items:            items,
			Min:              min,
			Redirect:         d["redirect"] == "1",
			Cache:            d["cache"] == "1",
			GmapKey:          d["gmap_key"],
			FuzzySearch:      d["fuzzy_search"] == "1",
			ExactMatchSearch: d["exact_match_search"] == "1",
			Color:            d["color"],
			ImageWidth:       width,
			ImageHeight:      height,
		},
	}
}

func (s Contact) Validate(ctx context.Context) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "validating a contact settings")

	if err := validators.V.Struct(s); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the shop", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	return nil
}

func (s ShopSettings) Validate(ctx context.Context) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "validating a shop settings")

	if err := validators.V.Struct(s); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the shop", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	return nil
}

func (s Contact) Save(ctx context.Context) error {
	l := slog.With(slog.String("name", s.Name))
	l.LogAttrs(ctx, slog.LevelInfo, "trying to save the contact shop")

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
		"banner_1", s.Banner1,
		"banner_2", s.Banner2,
		"banner_3", s.Banner3,
		"updated_at", now.Unix(),
	).Result()

	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot save the shop", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	return nil
}

func (s ShopSettings) Save(ctx context.Context) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "trying to save the shop")

	guest := "0"
	if s.Guest {
		guest = "1"
	}

	quantity := "0"
	if s.Quantity {
		quantity = "1"
	}

	new := "0"
	if s.New {
		new = "1"
	}

	redirect := "0"
	if s.Redirect {
		redirect = "1"
	}

	fuzzySearch := "0"
	if s.FuzzySearch {
		fuzzySearch = "1"
	}

	exactMatchSearch := "0"
	if s.ExactMatchSearch {
		exactMatchSearch = "1"
	}

	cache := "0"
	if s.Cache {
		cache = "1"
	}

	now := time.Now()
	_, err := db.Redis.HSet(context.Background(), "shop",
		"guest", guest,
		"quantity", quantity,
		"new", new,
		"items", fmt.Sprintf("%d", s.Items),
		"min", fmt.Sprintf("%d", s.Min),
		"redirect", redirect,
		"cache", cache,
		"gmap_key", s.GmapKey,
		"fuzzy_search", fuzzySearch,
		"exact_match_search", exactMatchSearch,
		"color", s.Color,
		"image_width", s.ImageWidth,
		"image_height", s.ImageHeight,
		"updated_at", now.Unix(),
	).Result()

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot save the shop", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	return nil
}
