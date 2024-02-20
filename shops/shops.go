package shops

import (
	"artisons/conf"
	"artisons/db"
	"artisons/validators"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type Contact struct {
	Name      string `validate:"required"`
	Address   string
	City      string
	Zipcode   string
	Phone     string `validate:"required"`
	Email     string `validate:"email"`
	Logo      string
	Banner1   string
	Banner2   string
	Banner3   string
	UpdatedAt time.Time
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
	Min float64

	ThrowsWhenPaymentFailed bool

	DeliveryFees float64

	DeliveryFreeFees float64

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

	updatedAt, err := strconv.ParseInt(d["updated_at"], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the updated at", slog.String("error", err.Error()), slog.String("updated_at", d["updated_at"]))
	}

	Data = Settings{
		Contact: Contact{
			Name:      d["name"],
			Address:   d["address"],
			City:      d["city"],
			Zipcode:   d["zipcode"],
			Phone:     d["phone"],
			Email:     d["email"],
			Logo:      d["logo"],
			Banner1:   d["banner_1"],
			Banner2:   d["banner_2"],
			Banner3:   d["banner_3"],
			UpdatedAt: time.Unix(updatedAt, 0),
		},
		ShopSettings: parseShopSettings(ctx, d),
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

func (s Contact) Save(ctx context.Context) (string, error) {
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
		return "", errors.New("something went wrong")
	}

	Data.Contact = s

	return "", nil
}

func parseShopSettings(ctx context.Context, data map[string]string) ShopSettings {
	var items int = 0
	if data["items"] != "" {
		i, err := strconv.ParseInt(data["items"], 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the items", slog.String("items", data["items"]), slog.String("error", err.Error()))
		} else {
			items = int(i)
		}
	}

	var err error
	var min float64 = 0
	if data["min"] != "" {
		min, err = strconv.ParseFloat(data["min"], 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the min", slog.String("min", data["min"]), slog.String("error", err.Error()))
		}
	}

	var width int = 0
	if data["image_width"] != "" {
		i, err := strconv.ParseInt(data["image_width"], 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the image_width", slog.String("image_width", data["image_width"]), slog.String("error", err.Error()))
		} else {
			width = int(i)
		}
	}

	var height int = 0
	if data["image_height"] != "" {
		i, err := strconv.ParseInt(data["image_height"], 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the image_height", slog.String("image_height", data["image_width"]), slog.String("error", err.Error()))
		} else {
			height = int(i)
		}
	}

	var deliveryFees float64 = 0
	if data["delivery_fees"] != "" {
		deliveryFees, err = strconv.ParseFloat(data["delivery_fees"], 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the min", slog.String("delivery_fees", data["delivery_fees"]), slog.String("error", err.Error()))
		}
	}

	var deliveryFreeFees float64 = 0
	if data["delivery_free_fees"] != "" {
		deliveryFreeFees, err = strconv.ParseFloat(data["delivery_free_fees"], 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the min", slog.String("delivery_free_fees", data["delivery_free_fees"]), slog.String("error", err.Error()))
		}
	}

	return ShopSettings{
		Guest:                   data["guest"] == "1",
		Quantity:                data["quantity"] == "1",
		New:                     data["new"] == "1",
		Items:                   items,
		Min:                     min,
		Redirect:                data["redirect"] == "1",
		Cache:                   data["cache"] == "1",
		GmapKey:                 data["gmap_key"],
		FuzzySearch:             data["fuzzy_search"] == "1",
		ExactMatchSearch:        data["exact_match_search"] == "1",
		Color:                   data["color"],
		ImageWidth:              width,
		ImageHeight:             height,
		DeliveryFees:            deliveryFees,
		DeliveryFreeFees:        deliveryFreeFees,
		ThrowsWhenPaymentFailed: data["throws_when_payment_failed"] == "1",
	}
}

func (s ShopSettings) Save(ctx context.Context) (string, error) {
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

	throwsWhenPaymentFailed := "0"
	if s.ThrowsWhenPaymentFailed {
		throwsWhenPaymentFailed = "1"
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
		"min", fmt.Sprintf("%2f", s.Min),
		"redirect", redirect,
		"cache", cache,
		"gmap_key", s.GmapKey,
		"fuzzy_search", fuzzySearch,
		"exact_match_search", exactMatchSearch,
		"color", s.Color,
		"image_width", s.ImageWidth,
		"image_height", s.ImageHeight,
		"delivery_fees", s.DeliveryFees,
		"delivery_free_fees", s.DeliveryFreeFees,
		"throws_when_payment_failed", throwsWhenPaymentFailed,
		"updated_at", now.Unix(),
	).Result()

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot save the shop", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	Data.ShopSettings = s

	return "", nil
}

func Deliveries(ctx context.Context) ([]string, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "getting deliveries")

	del, err := db.Redis.ZRange(ctx, "deliveries", 0, 999).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the deliveries", slog.String("error", err.Error()))
		return []string{}, errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "got deliveries", slog.Int("length", len(del)))

	return del, nil
}

func Payments(ctx context.Context) ([]string, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "getting payments")

	pay, err := db.Redis.ZRange(ctx, "payments", 0, 999).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the payments", slog.String("error", err.Error()))
		return []string{}, errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "got payments", slog.Int("length", len(pay)))

	return pay, nil
}

func Pay(ctx context.Context, oid string, payment string) (string, error) {
	l := slog.With(slog.String("order", oid), slog.String("payment", payment))

	switch payment {
	case "cash":
		{
			l.LogAttrs(ctx, slog.LevelInfo, "payment success")
			return "", nil
		}
	default:
		{
			l.LogAttrs(ctx, slog.LevelError, "the payment method does not exist")
			return "", errors.New("you are not authorized to process this request")
		}
	}

}

// IsValidDelivery returns true if the delivery
// is valid. The values can be "collect" or "home".
// The "collect" value can be used only if it's allowed
// in the settings.
func IsValidDelivery(ctx context.Context, d string) bool {
	l := slog.With(slog.String("delivery", d))
	l.LogAttrs(ctx, slog.LevelInfo, "checking delivery validity")

	del, err := Deliveries(ctx)
	if err != nil {
		return false
	}

	if err := validators.V.Var(d, "oneof="+strings.Join(del, " ")); err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate  the delivery", slog.String("error", err.Error()))
		return false
	}

	if d == "home" && !conf.HasHomeDelivery {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot continue while the home is not activated")
		return false
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the delivery is valid")

	return true
}

// IsValidPayment returns true if the payment
// is valid. The values can be "card", "cash", "bitcoin" or "wire".
// The payments can be enablee or disabled in the settings.
func IsValidPayment(ctx context.Context, p string) bool {
	l := slog.With(slog.String("payment", p))
	l.LogAttrs(ctx, slog.LevelInfo, "checking payment validity")

	pay, err := Payments(ctx)
	if err != nil {
		return false
	}

	if err := validators.V.Var(p, "oneof="+strings.Join(pay, " ")); err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the payment", slog.String("error", err.Error()))
		return false
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the payment is valid")

	return true
}

func DeliveryFreeFees(ctx context.Context) (float64, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "retrieving shop delivery free fees")

	d, err := db.Redis.HGet(ctx, "shop", "delivery_free_fees").Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get shop delivery free fees info", slog.String("error", err.Error()))
		return 0, errors.New("something went wrong")
	}

	val, err := strconv.ParseFloat(d, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse shop delivery free fees info", slog.String("error", err.Error()))
		return 0, errors.New("something went wrong")
	}

	return val, nil
}

func DeliveryFees(ctx context.Context) (float64, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "retrieving shop delivery fees")

	d, err := db.Redis.HGet(ctx, "shop", "delivery_fees").Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get shop delivery free fees info", slog.String("error", err.Error()))
		return 0, errors.New("something went wrong")
	}

	val, err := strconv.ParseFloat(d, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse shop delivery free fees info", slog.String("error", err.Error()))
		return 0, errors.New("something went wrong")
	}

	return val, nil
}
