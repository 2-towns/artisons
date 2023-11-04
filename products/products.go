// Package products provide everything around products
package products

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/locales"
	"log"
	"log/slog"
	"regexp"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Product is the product representation in the application
type Product struct {
	PID         string  `redis:"pid"` // ID is an unique identifier
	Title       string  `redis:"title"`
	Description string  `redis:"description"`
	Price       float32 `redis:"price"`
	Slug        string  `redis:"slug"`
	MID         string  `redis:"mid"`
	Sku         string  `redis:"sku"`
	Currency    string  `redis:"currency"`
	Quantity    int     `redis:"quantity"`
	// Images length
	Length int     `redis:"length"`
	Status string  `redis:"status"`
	Weight float32 `redis:"weight"`

	Images []string
	Tags   []string
	Links  []string          // Links contains the linked product IDs
	Meta   map[string]string // Meta contains the product options.
}

const (
	Online  = "online"  // Make th product available in the application
	Offline = "offline" // Hide th product  in the application
)

// ImageExtensions is the allowed extensions in the application
const ImageExtensions = "jpg jpeg png"

// GetImagePath returns the imgproxy path for a file
// Later on, the method should be improve to generate subfolders path,
// if the products are more than the unix file limit
func ImagePath(pid string, index int) (string, string) {
	folder := fmt.Sprintf("%s/%s", conf.ImgProxyPath, pid)
	return folder, fmt.Sprintf("%s/%d", folder, index)
}

// Available return true if all the product ids are availables
func Availables(c context.Context, pids []string) bool {
	l := slog.With(slog.Any("pids", pids))
	l.LogAttrs(c, slog.LevelInfo, "checking the pids availability")

	ctx := context.Background()
	pipe := db.Redis.Pipeline()
	for _, pid := range pids {
		pipe.HGet(ctx, "product:"+pid, "status")
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the product ids", slog.String("error", err.Error()))
		return false
	}

	for _, cmd := range cmds {
		status := cmd.(*redis.StringCmd).Val()
		if status != "online" {
			l.LogAttrs(c, slog.LevelInfo, "cannot get the product while it is not available", slog.String("id", cmd.Args()[1].(string)))
			return false
		}
	}

	l.LogAttrs(c, slog.LevelInfo, "the pids are available")

	return true
}

// Available return true if the product is available
func Available(c context.Context, pid string) bool {
	l := slog.With(slog.String("pid", pid))
	l.LogAttrs(c, slog.LevelInfo, "checking the pid availability")

	ctx := context.Background()

	if exists, err := db.Redis.Exists(ctx, "product:"+pid).Result(); exists == 0 || err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the product")
		return false
	}

	status, err := db.Redis.HGet(ctx, "product:"+pid, "status").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot find the product", slog.String("error", err.Error()))
		return false
	}

	l.LogAttrs(c, slog.LevelInfo, "got the product status", slog.String("availability", "status"))

	return status == "online"
}

func parse(c context.Context, data, options map[string]string, tags, links []string) (Product, error) {
	slog.LogAttrs(c, slog.LevelInfo, "parsing the product data")

	l, err := strconv.ParseInt(data["length"], 10, 8)
	if err != nil {
		slog.Error("cannot parse the product length", slog.String("length", data["length"]))
		return Product{}, err
	}
	length := int(l)

	price, err := strconv.ParseFloat(data["price"], 32)
	if err != nil {
		slog.Error("cannot parse the product price", slog.String("price", data["price"]))
		return Product{}, err
	}

	quantity, err := strconv.ParseInt(data["quantity"], 10, 8)
	if err != nil {
		slog.Error("cannot parse the product quantity", slog.String("quantity", data["quantity"]))
		return Product{}, err
	}

	var weight float32

	if data["weight"] != "" {
		v, err := strconv.ParseFloat(data["weight"], 32)
		if err != nil {
			slog.Error("cannot parse the product weight", slog.String("weight", data["weight"]))
			return Product{}, err
		}

		weight = float32(v)
	}

	images := []string{}
	for i := 0; i < length; i++ {
		_, image := ImagePath(data["id"], i)
		images = append(images, image)
	}

	slog.LogAttrs(c, slog.LevelInfo, "product parsed successfully")

	return Product{
		PID:         data["pid"],
		Title:       data["title"],
		Description: data["description"],
		Price:       float32(price),
		Slug:        data["slug"],
		MID:         data["mid"],
		Sku:         data["sku"],
		Currency:    data["currency"],
		Quantity:    int(quantity),
		Weight:      float32(weight),
		Status:      data["status"],
		Links:       links,
		Tags:        tags,
		Meta:        options,
		Length:      length,
		Images:      images,
	}, nil
}

func (p Product) Validate(c context.Context) error {
	slog.LogAttrs(c, slog.LevelInfo, "validating a product")
	log.Println(c.Value(locales.ContextKey))
	var printer = message.NewPrinter(c.Value(locales.ContextKey).(language.Tag))

	if p.Sku == "" {
		slog.Info("cannot parse the empty sku")
		log.Println(printer.Sprintf("input_required", "sku"))
		return errors.New(printer.Sprintf("input_required", "sku"))
	}

	isValid := regexp.MustCompile(`^[0-9a-z]+$`).MatchString
	if !isValid(p.Sku) {
		slog.Info("cannot parse the empty sku", slog.String("sku", p.Sku))
		return errors.New(printer.Sprintf("input_validation", "sku"))
	}

	if p.Title == "" {
		slog.Info("cannot parse the empty title")
		return errors.New(printer.Sprintf("input_required", "title"))
	}

	if !conf.IsCurrencySupported(p.Currency) {
		slog.Info("cannot use an unsupported currency", slog.String("currency", p.Currency))
		return errors.New(printer.Sprintf("input_validation", "currency"))
	}

	if p.Status != "online" && p.Status != "offline" {
		slog.Error("cannot use an unsupported status", slog.String("status", p.Status))
		return errors.New(printer.Sprintf("input_validation", "status"))
	}

	if p.Description == "" {
		slog.Info("cannot parse the empty description")
		return errors.New(printer.Sprintf("input_required", "description"))
	}

	if p.Length <= 0 {
		slog.Info("cannot parse the empty length")
		return errors.New(printer.Sprintf("input_required", "length"))
	}

	return nil
}

func (p Product) Save(ctx context.Context) error {
	if p.PID == "" {
		slog.Error("cannot continue with empty pid")
		return errors.New("product_pid_required")
	}

	l := slog.With(slog.String("pid", p.PID))
	l.Info("storing the product")

	key := "product:" + p.PID
	lkey := key + ":links"
	tkey := key + ":tags"
	okey := key + ":options"
	now := time.Now()

	_, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, lkey, okey, tkey)
		rdb.HSet(ctx, key,
			"pid", p.PID,
			"sku", p.Sku,
			"title", p.Title,
			"description", p.Description,
			"length", p.Length,
			"currency", p.Currency,
			"price", p.Price,
			"quantity", p.Quantity,
			"status", p.Status,
			"weight", p.Weight,
			"mid", p.MID,
			"updated_at", time.Now().Format(time.RFC3339),
		)
		rdb.ZAdd(ctx, "products:"+p.MID, redis.Z{
			Score:  float64(now.Unix()),
			Member: p.PID,
		})

		if len(p.Links) > 0 {
			rdb.SAdd(ctx, lkey, p.Links)
		}

		if len(p.Tags) > 0 {
			rdb.SAdd(ctx, tkey, p.Tags)
		}

		if len(p.Meta) > 0 {
			for k, v := range p.Meta {
				rdb.HSet(ctx, okey, k, v)
			}
		}

		return nil
	})

	if err != nil {
		slog.Error("cannot store the product", slog.String("error", err.Error()))
	} else {
		l.Info("product stored successfully")
	}

	return err
}

// Find looks for a product by its product id
func Find(c context.Context, pid string) (Product, error) {
	l := slog.With(slog.String("pid", pid))
	l.LogAttrs(c, slog.LevelInfo, "looking for product")

	ctx := context.Background()

	if exists, err := db.Redis.Exists(ctx, "product:"+pid).Result(); exists == 0 || err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the product")
		return Product{}, errors.New("product_not_found")
	}

	data, err := db.Redis.HGetAll(ctx, "product:"+pid).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot find the product", slog.String("error", err.Error()))
		return Product{}, err
	}

	tags, err := db.Redis.SMembers(ctx, "product:"+pid+":tags").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot find the product tags", slog.String("error", err.Error()))
		return Product{}, err
	}

	links, err := db.Redis.SMembers(ctx, "product:"+pid+":links").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot find the product links", slog.String("error", err.Error()))
		return Product{}, err
	}

	options, err := db.Redis.HGetAll(ctx, "product:"+pid+":options").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot find the product options", slog.String("error", err.Error()))
		return Product{}, err
	}

	p, err := parse(c, data, options, tags, links)

	if err != nil {
		l.LogAttrs(c, slog.LevelInfo, "the product is found", slog.String("sku", p.Sku))
	}

	return p, err
}
