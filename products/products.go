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
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Product is the product representation in the application
type Product struct {
	ID          string  `redis:"id"` // ID is an unique identifier
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

type Meta = map[string]string

type Query struct {
	Keywords string
	PriceMin float32
	PriceMax float32
	Tags     []string
	Meta     map[string]string
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
	l := slog.With(slog.Any("ids", pids))
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
	l := slog.With(slog.String("id", pid))
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

func parse(c context.Context, data map[string]string) (Product, error) {
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
		ID:          data["id"],
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
		Tags:        strings.Split(data["status"], ";"),
		Links:       strings.Split(data["links"], ";"),
		Meta:        UnSerializeMeta(c, data["meta"], ";"),
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
	if p.ID == "" {
		slog.Error("cannot continue with empty pid")
		return errors.New("product_pid_required")
	}

	l := slog.With(slog.String("id", p.ID))
	l.Info("storing the product")

	key := "product:" + p.ID
	now := time.Now()

	_, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, key,
			"id", p.ID,
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
			"tags", strings.Join(p.Tags, ";"),
			"links", strings.Join(p.Links, ";"),
			"meta", SerializeMeta(ctx, p.Meta, ";"),
			"created_at", time.Now().Format(time.RFC3339),
			"updated_at", time.Now().Format(time.RFC3339),
		)
		rdb.ZAdd(ctx, "products:"+p.MID, redis.Z{
			Score:  float64(now.Unix()),
			Member: p.ID,
		})

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
	l := slog.With(slog.String("id", pid))
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

	p, err := parse(c, data)

	if err != nil {
		l.LogAttrs(c, slog.LevelInfo, "the product is found", slog.String("sku", p.Sku))
	}

	return p, err
}

func or(qs string, s string) string {
	if qs == "" {
		return s
	}

	return qs + " OR " + s
}

func convertMap(m map[interface{}]interface{}) map[string]string {
	v := map[string]string{}

	for key, value := range m {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		v[strKey] = strValue
	}

	return v
}

func Search(c context.Context, q Query) ([]Product, error) {
	slog.LogAttrs(c, slog.LevelInfo, "searching products")

	qs := "@status:{online} "

	if q.Keywords != "" {
		qs += fmt.Sprintf("(@title:*%s*)|(@description:*%s*)|(@sku:{%s})", q.Keywords, q.Keywords, q.Keywords)
	}

	var priceMin interface{} = "-inf"
	var priceMax interface{} = "+inf"

	if q.PriceMin > 0 {
		priceMin = q.PriceMin
	}

	if q.PriceMax > 0 {
		priceMax = q.PriceMax
	}

	if priceMin != "-inf" || priceMax != "+inf" {
		qs += fmt.Sprintf("@price:[%v %v]", priceMin, priceMax)
	}

	if len(q.Tags) > 0 {
		qs += fmt.Sprintf("@tags:{%s}", strings.Join(q.Tags, " | "))
	}

	if len(q.Meta) > 0 {
		s := SerializeMeta(c, q.Meta, " | ")
		qs += fmt.Sprintf("@meta:{%s}", s)
	}

	slog.LogAttrs(c, slog.LevelInfo, "preparing redis request", slog.String("query", qs))

	ctx := context.Background()
	cmds, err := db.Redis.Do(
		ctx,
		"FT.SEARCH",
		db.SearchIdx,
		qs,
	).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot run the search query", slog.String("query", qs), slog.String("error", err.Error()))
		log.Fatal()
	}

	res := cmds.(map[interface{}]interface{})
	slog.LogAttrs(c, slog.LevelInfo, "search done", slog.Int64("results", res["total_results"].(int64)))

	results := res["results"].([]interface{})
	products := []Product{}

	for _, value := range results {
		m := value.(map[interface{}]interface{})
		attributes := m["extra_attributes"].(map[interface{}]interface{})
		data := convertMap(attributes)

		product, err := parse(c, data)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the product", slog.Any("product", data), slog.String("error", err.Error()))
			continue
		}

		products = append(products, product)
	}

	return products, nil
}

// SerializeMeta transforms a meta map to a string representation.
// The values are separated by ";".
// Example: map["color"]"blue" => color_blue
func SerializeMeta(c context.Context, m map[string]string, sep string) string {
	slog.LogAttrs(c, slog.LevelInfo, "serializing the product meta", slog.Any("meta", m))

	s := ""
	for key, value := range m {
		if s != "" {
			s += sep
		}

		s += fmt.Sprintf("%s_%s", key, value)
	}

	slog.LogAttrs(c, slog.LevelInfo, "serialize done successfully", slog.String("serialized", s))

	return s
}

// UnSerializeMeta transform the meta serialized to a map.
// The values are separated by ";".
// Example: color_blue => map["color"]"blue"
func UnSerializeMeta(c context.Context, s, sep string) map[string]string {
	slog.LogAttrs(c, slog.LevelInfo, "unserializing the product meta", slog.String("serialized", s))
	values := strings.Split(s, sep)
	meta := map[string]string{}

	for _, value := range values {
		parts := strings.Split(value, "_")

		if len(parts) != 2 {
			slog.LogAttrs(c, slog.LevelError, "cannot unserialize the product meta", slog.String("serialized", s))
			continue
		}

		meta[parts[0]] = parts[1]
	}

	slog.LogAttrs(c, slog.LevelInfo, "unserialize done successfully")

	return meta
}
