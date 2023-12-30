// Package products provide everything around products
package products

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"gifthub/string/stringutil"
	"gifthub/tracking"
	"gifthub/users"
	"gifthub/validators"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

// Product is the product representation in the application
type Product struct {
	ID          string  `redis:"id"` // ID is an unique identifier
	Title       string  `redis:"title" validate:"required"`
	Description string  `redis:"description" validate:"required"`
	Price       float64 `redis:"price" validate:"required"`
	// The percent discount
	Discount float64 `redis:"discount"`
	Slug     string  `redis:"slug"`
	MID      string  `redis:"mid"`
	Sku      string  `redis:"sku" validate:"omitempty,alphanum"`
	Currency string  `redis:"currency"`
	Quantity int     `redis:"quantity" validate:"required"`
	Status   string  `redis:"status" validate:"oneof=online offline"`
	Weight   float64 `redis:"weight"`

	Image1 string
	Image2 string
	Image3 string
	Image4 string
	Tags   []string

	// todo manage the product links
	Links []string // Links contains the linked product IDs

	// todo manage the product links
	Meta map[string]string // Meta contains the product options.

	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchResults struct {
	Total    int64
	Products []Product
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

// ImageURL returns the imgproxy URL
func ImageURL(id string) string {
	if id == "" {
		return ""
	}

	folder := fmt.Sprintf("%s/%s", conf.ImgProxyURL, id)
	return folder
}

// GetImagePath returns the imgproxy path for a file
// Later on, the method should be improve to generate subfolders path,
// if the products are more than the unix file limit/
// The path should be on a different server
func ImagePath(id string) string {
	if id == "" {
		return ""
	}

	folder := fmt.Sprintf("%s/%s", conf.ImgProxyPath, id)
	return folder
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
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(c, slog.LevelError, "cannot get the status", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

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
	price, err := strconv.ParseFloat(data["price"], 32)
	if err != nil {
		slog.Error("cannot parse the product price", slog.String("price", data["price"]))
		return Product{}, errors.New("input_price_invalid")
	}

	quantity, err := strconv.ParseInt(data["quantity"], 10, 32)
	if err != nil {
		slog.Error("cannot parse the product quantity", slog.String("quantity", data["quantity"]))
		return Product{}, errors.New("input_quantity_invalid")
	}

	var weight float64

	if data["weight"] != "" {
		v, err := strconv.ParseFloat(data["weight"], 32)
		if err != nil {
			slog.Error("cannot parse the product weight", slog.String("weight", data["weight"]))
			return Product{}, errors.New("input_weight_invalid")
		}

		weight = v
	}

	if data["discount"] != "" {
		v, err := strconv.ParseFloat(data["discount"], 32)
		if err != nil {
			slog.Error("cannot parse the product discount", slog.String("discount", data["discount"]))
			return Product{}, errors.New("input_discount_invalid")
		}

		weight = v
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		slog.Error("cannot parse the product created at", slog.String("created_at", data["created_at"]))
		return Product{}, errors.New("input_created_at_invalid")
	}

	updatedAt, err := strconv.ParseInt(data["updated_at"], 10, 64)
	if err != nil {
		slog.Error("cannot parse the product updatede at", slog.String("updated_at", data["updated_at"]))
		return Product{}, errors.New("input_updated_at_invalid")
	}

	return Product{
		ID:          data["id"],
		Title:       db.Unescape(data["title"]),
		Description: db.Unescape(data["description"]),
		Price:       price,
		Slug:        db.Unescape(data["slug"]),
		MID:         data["mid"],
		Sku:         db.Unescape(data["sku"]),
		Currency:    data["currency"],
		Quantity:    int(quantity),
		Weight:      weight,
		Status:      data["status"],
		Tags:        strings.Split(db.Unescape(data["tags"]), ";"),
		Links:       strings.Split(db.Unescape(data["links"]), ";"),
		Meta:        UnSerializeMeta(c, db.Unescape(data["meta"]), ";"),
		Image1:      ImageURL(data["image_1"]),
		Image2:      ImageURL(data["image_2"]),
		Image3:      ImageURL(data["image_3"]),
		Image4:      ImageURL(data["image_4"]),
		CreatedAt:   time.Unix(createdAt, 0),
		UpdatedAt:   time.Unix(updatedAt, 0),
	}, nil
}

func (p Product) Validate(c context.Context) error {
	slog.LogAttrs(c, slog.LevelInfo, "validating a product")

	if err := validators.V.Struct(p); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot validate the product", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input_%s_invalid", low)
	}

	if !conf.IsCurrencySupported(p.Currency) {
		slog.Info("cannot use an unsupported currency", slog.String("currency", p.Currency))
		return errors.New("input_currency_invalid")
	}

	return nil
}

// Save a product into redis.
// The keys are :
// product:pid => the product data
func (p Product) Save(ctx context.Context) error {
	if p.ID == "" {
		slog.LogAttrs(ctx, slog.LevelError, "cannot continue with empty pid")
		return errors.New("input_pid_required")
	}

	l := slog.With(slog.String("id", p.ID))

	// score, err := db.Redis.ZScore(ctx, "products", p.ID).Result()
	// if err != nil {
	// 	l.LogAttrs(ctx, slog.LevelError, "cannot verify product existence", slog.String("error", err.Error()))
	// }

	// if score == 0 {
	// 	l.LogAttrs(ctx, slog.LevelInfo, "the product is new")
	// }

	key := "product:" + p.ID
	title := db.Escape(p.Title)

	if _, err := db.Redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key,
			"id", p.ID,
			"sku", db.Escape(p.Sku),
			"title", title,
			"slug", stringutil.Slugify(title),
			"description", db.Escape(p.Title),
			"currency", p.Currency,
			"price", p.Price,
			"quantity", p.Quantity,
			"status", p.Status,
			"weight", p.Weight,
			"mid", p.MID,
			"tags", db.Escape(strings.Join(p.Tags, ";")),
			// "links", db.Escape(strings.Join(p.Links, ";")),
			// "meta", db.Escape(SerializeMeta(ctx, p.Meta, ";")),
			"created_at", time.Now().Unix(),
			"updated_at", time.Now().Unix(),
		)

		if p.Image1 == "-" {
			pipe.HSet(ctx, key, "image_1", "")
		} else if p.Image1 != "" {
			pipe.HSet(ctx, key, "image_1", p.Image1)
		}

		if p.Image2 == "-" {
			pipe.HSet(ctx, key, "image_2", "")
		} else if p.Image2 != "" {
			pipe.HSet(ctx, key, "image_2", p.Image2)
		}

		if p.Image3 == "-" {
			pipe.HSet(ctx, key, "image_3", "")
		} else if p.Image3 != "" {
			pipe.HSet(ctx, key, "image_3", p.Image3)
		}

		if p.Image4 == "-" {
			pipe.HSet(ctx, key, "image_4", "")
		} else if p.Image4 != "" {
			pipe.HSet(ctx, key, "image_4", p.Image4)
		}

		// if score == 0 {
		// 	pipe.ZAdd(ctx, "products", redis.Z{
		// 		Score:  float64(time.Now().Unix()),
		// 		Member: p.ID,
		// 	})
		// }

		return nil
	}); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot store the product", slog.String("error", err.Error()))
		return err
	}

	return nil
}

// Find looks for a product by its product id
func Find(c context.Context, pid string) (Product, error) {
	l := slog.With(slog.String("id", pid))
	l.LogAttrs(c, slog.LevelInfo, "looking for product")

	if pid == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate empty product id")
		return Product{}, errors.New("input_id_required")
	}

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

func Search(c context.Context, q Query, offset, num int) (SearchResults, error) {
	slog.LogAttrs(c, slog.LevelInfo, "searching products")

	qs := "@status:{online} "

	if q.Keywords != "" {
		k := db.Escape(q.Keywords)
		qs += fmt.Sprintf("(@title:'*%s*')|(@description:'*%s*')|(@sku:'{%s}')|(@id:'{%s})'", k, k, k, k)
	}

	var priceMin interface{} = "-inf"
	var priceMax interface{} = "+inf"
	priceMinRep := "%v"
	priceMaxRep := "%v"

	if q.PriceMin > 0 {
		priceMinRep = "%f"
		priceMin = q.PriceMin
	}

	if q.PriceMax > 0 {
		priceMaxRep = "%f"
		priceMax = q.PriceMax
	}

	if priceMin != "-inf" || priceMax != "+inf" {
		qs += fmt.Sprintf("@price:["+priceMinRep+" "+priceMaxRep+"]", priceMin, priceMax)
	}

	if len(q.Tags) > 0 {
		t := db.Escape(strings.Join(q.Tags, " | "))
		qs += fmt.Sprintf("@tags:{%s}", t)
	}

	if len(q.Meta) > 0 {
		s := db.Escape(SerializeMeta(c, q.Meta, " | "))
		qs += fmt.Sprintf("@meta:{%s}", s)
	}

	slog.LogAttrs(c, slog.LevelInfo, "preparing redis request", slog.String("query", qs))

	ctx := context.Background()
	cmds, err := db.Redis.Do(
		ctx,
		"FT.SEARCH",
		db.ProductIdx,
		qs,
		"LIMIT",
		fmt.Sprintf("%d", offset),
		fmt.Sprintf("%d", num),
		"SORTBY",
		"updated_at",
		"desc",
	).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot run the search query", slog.String("query", qs), slog.String("error", err.Error()))
		return SearchResults{}, err
	}

	res := cmds.(map[interface{}]interface{})
	total := res["total_results"].(int64)

	slog.LogAttrs(c, slog.LevelInfo, "search done", slog.Int64("results", total))

	results := res["results"].([]interface{})
	products := []Product{}

	for _, value := range results {
		m := value.(map[interface{}]interface{})
		attributes := m["extra_attributes"].(map[interface{}]interface{})
		data := db.ConvertMap(attributes)

		product, err := parse(c, data)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the product", slog.Any("product", data), slog.String("error", err.Error()))
			continue
		}

		products = append(products, product)
	}

	user, ok := c.Value(contexts.User).(users.User)
	if !ok || user.Role != "admin" {
		tra := map[string]string{
			"query": fmt.Sprintf("'%s'", qs),
		}

		go tracking.Log(c, "product_search", tra)
	}

	return SearchResults{
		Total:    total,
		Products: products,
	}, nil
}

func Count(c context.Context) (int64, error) {
	slog.LogAttrs(c, slog.LevelInfo, "counting products")

	count, err := db.Redis.ZCount(c, "products", "-inf", "+inf").Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot find the product", slog.String("error", err.Error()))
		return 0, err
	}

	return count, nil
}

// func List(c context.Context, offset int) ([]Product, error) {
// l := slog.With(slog.Int("page", offset))
// l.LogAttrs(c, slog.LevelInfo, "listing products")

// ctx := context.Background()
// start := int64(offset * conf.ItemsPerPage)
// end := int64(start + conf.ItemsPerPage)

// pids, err := db.Redis.ZRevRange(ctx, "products", start, end).Result()
// if err != nil {
// 	l.LogAttrs(c, slog.LevelError, "cannot find the product", slog.String("error", err.Error()))
// 	return []Product{}, err
// }

// pipe := db.Redis.Pipeline()
// for _, pid := range pids {
// 	pipe.HGetAll(ctx, "product:"+pid)
// }

// cmds, err := pipe.Exec(ctx)
// if err != nil {
// 	l.LogAttrs(c, slog.LevelError, "cannot get the product ids", slog.String("error", err.Error()))
// 	return []Product{}, err
// }

// pds := []Product{}

// for _, cmd := range cmds {
// 	key := fmt.Sprintf("%s", cmd.Args()[1])

// 	if cmd.Err() != nil {
// 		slog.LogAttrs(c, slog.LevelError, "cannot get the product", slog.String("key", key), slog.String("error", err.Error()))
// 		continue
// 	}

// 	m := cmd.(*redis.MapStringStringCmd).Val()
// 	p, err := parse(ctx, m)
// 	if err != nil {
// 		l.LogAttrs(c, slog.LevelError, "cannot get the product ids", slog.String("error", err.Error()))
// 		continue
// 	}

// 	pds = append(pds, p)
// }

// return pds, nil
// }

func (p Product) URL() string {
	return conf.WebsiteURL + "/" + p.ID + "-" + p.Slug + ".html"
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
	if s == "" {
		return map[string]string{}
	}

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

	return meta
}

func Delete(ctx context.Context, pid string) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "deleting product", slog.Any("pid", pid))

	if pid == "" {
		slog.LogAttrs(ctx, slog.LevelInfo, "the pid cannot be empty")
		return errors.New("input_id_invalid")
	}

	if _, err := db.Redis.Del(ctx, "product:"+pid).Result(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot delete product", slog.String("string", err.Error()))
		return err
	}

	return nil
}
