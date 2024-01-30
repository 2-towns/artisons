// Package products provide everything around products
package products

import (
	"artisons/conf"
	"artisons/db"
	"artisons/string/stringutil"
	"artisons/validators"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path"
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
	Slug     string  `redis:"slug" validate:"required"`
	MID      string  `redis:"mid"`
	Sku      string  `redis:"sku" validate:"omitempty,alphanum"`
	Quantity int     `redis:"quantity" validate:"required"`
	Status   string  `redis:"status" validate:"oneof=online offline"`
	Weight   float64 `redis:"weight"`

	Image1 string
	Image2 string
	Image3 string
	Image4 string
	Tags   []string

	// Links []string // Links contains the linked product IDs

	Meta map[string][]string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchResults struct {
	Total    int
	Products []Product
}

type Meta = map[string]string

type Query struct {
	Keywords string
	PriceMin float32
	PriceMax float32
	Tags     []string
	Meta     map[string][]string
	Slug     string
}

const (
	Online  = "online"  // Make th product available in the application
	Offline = "offline" // Hide th product  in the application
)

// ImageExtensions is the allowed extensions in the application
const ImageExtensions = "jpg jpeg png"

// GetImagePath returns the imgproxy path for a file
// Later on, the method should be improve to generate subfolders path,
// if the products are more than the unix file limit/
// The path should be on a different server
func ImagePath(id string) string {
	if id == "" {
		return ""
	}

	return path.Join(conf.ImgProxy.Path, id)
}

// Available return true if all the product ids are availables
func Availables(ctx context.Context, pids []string) bool {
	l := slog.With(slog.Any("ids", pids))
	l.LogAttrs(ctx, slog.LevelInfo, "checking the pids availability")

	pipe := db.Redis.Pipeline()
	for _, pid := range pids {
		pipe.HGet(ctx, "product:"+pid, "status")
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the product ids", slog.String("error", err.Error()))
		return false
	}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the status", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		status := cmd.(*redis.StringCmd).Val()
		if status != "online" {
			l.LogAttrs(ctx, slog.LevelInfo, "cannot get the product while it is not available", slog.String("id", cmd.Args()[1].(string)))
			return false
		}
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the pids are available")

	return true
}

// Available return true if the product is available
func Available(ctx context.Context, pid string) bool {
	l := slog.With(slog.String("id", pid))
	l.LogAttrs(ctx, slog.LevelInfo, "checking the pid availability")

	if exists, err := db.Redis.Exists(ctx, "product:"+pid).Result(); exists == 0 || err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot find the product")
		return false
	}

	status, err := db.Redis.HGet(ctx, "product:"+pid, "status").Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot find the product", slog.String("error", err.Error()))
		return false
	}

	l.LogAttrs(ctx, slog.LevelInfo, "got the product status", slog.String("availability", "status"))

	return status == "online"
}

func parse(ctx context.Context, data map[string]string) (Product, error) {
	price, err := strconv.ParseFloat(data["price"], 32)
	if err != nil {
		slog.Error("cannot parse the product price", slog.String("price", data["price"]))
		return Product{}, errors.New("input:price")
	}

	quantity, err := strconv.ParseInt(data["quantity"], 10, 32)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the product quantity", slog.String("quantity", data["quantity"]))
		return Product{}, errors.New("input:quantity")
	}

	var weight float64

	if data["weight"] != "" {
		v, err := strconv.ParseFloat(data["weight"], 32)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the product weight", slog.String("weight", data["weight"]))
			return Product{}, errors.New("input:weight")
		}

		weight = v
	}

	if data["discount"] != "" {
		v, err := strconv.ParseFloat(data["discount"], 32)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the product discount", slog.String("discount", data["discount"]))
			return Product{}, errors.New("input:discount")
		}

		weight = v
	}

	updatedAt, err := strconv.ParseInt(data["updated_at"], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the product updated at", slog.String("updated_at", data["updated_at"]))
		return Product{}, errors.New("input:updated_at")
	}

	return Product{
		ID:          data["id"],
		Title:       db.Unescape(data["title"]),
		Description: db.Unescape(data["description"]),
		Price:       price,
		Slug:        db.Unescape(data["slug"]),
		MID:         data["mid"],
		Sku:         db.Unescape(data["sku"]),
		Quantity:    int(quantity),
		Weight:      weight,
		Status:      data["status"],
		Tags:        strings.Split(db.Unescape(data["tags"]), ";"),
		Meta:        UnSerializeMeta(ctx, db.Unescape(data["meta"])),
		Image1:      data["image_1"],
		Image2:      data["image_2"],
		Image3:      data["image_3"],
		Image4:      data["image_4"],
		UpdatedAt:   time.Unix(updatedAt, 0),
	}, nil
}

func (p Product) Validate(ctx context.Context) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "validating a product")

	if err := validators.V.Struct(p); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the product", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	return nil
}

// Save a product into redis.
// The keys are :
// product:pid => the product data
func (p Product) Save(ctx context.Context) (string, error) {
	if p.ID == "" {
		pid, err := stringutil.Random()

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot generated the product id", slog.String("error", err.Error()))
			return "", errors.New("something went wrong")
		}

		p.ID = pid
	}

	l := slog.With(slog.String("id", p.ID))

	key := "product:" + p.ID
	title := db.Escape(p.Title)
	now := time.Now().Unix()

	var values []interface{}
	values = append(values,
		"sku", db.Escape(p.Sku),
		"title", title,
		"slug", db.Escape(p.Slug),
		"description", db.Escape(p.Title),
		"price", p.Price,
		"quantity", p.Quantity,
		"status", p.Status,
		"weight", p.Weight,
		"mid", p.MID,
		"tags", db.Escape(strings.Join(p.Tags, ";")),
		// "links", db.Escape(strings.Join(p.Links, ";")),
		"meta", db.Escape(SerializeMeta(ctx, p.Meta)),
		"updated_at", now,
	)

	if p.Image1 == "-" {
		values = append(values, "image_1", "")
	} else if p.Image1 != "" {
		values = append(values, "image_1", p.Image1)
	}

	if p.Image2 == "-" {
		values = append(values, "image_2", "")
	} else if p.Image2 != "" {
		values = append(values, "image_2", p.Image2)
	}

	if p.Image3 == "-" {
		values = append(values, "image_3", "")
	} else if p.Image3 != "" {
		values = append(values, "image_3", p.Image3)
	}

	if p.Image4 == "-" {
		values = append(values, "image_4", "")
	} else if p.Image4 != "" {
		values = append(values, "image_4", p.Image4)
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, key, values)
		rdb.HSetNX(ctx, key, "created_at", now)
		rdb.HSetNX(ctx, key, "id", p.ID)
		rdb.HSetNX(ctx, key, "type", "product")

		return nil
	}); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot store the product", slog.String("error", err.Error()))
		return "", err
	}

	return p.ID, nil
}

// Find looks for a product by its product id
func Find(ctx context.Context, pid string) (Product, error) {
	l := slog.With(slog.String("id", pid))
	l.LogAttrs(ctx, slog.LevelInfo, "looking for product")

	if pid == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate empty product id")
		return Product{}, errors.New("input:id")
	}

	if exists, err := db.Redis.Exists(ctx, "product:"+pid).Result(); exists == 0 || err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot find the product")
		return Product{}, errors.New("oops the data is not found")
	}

	data, err := db.Redis.HGetAll(ctx, "product:"+pid).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot find the product", slog.String("error", err.Error()))
		return Product{}, err
	}

	p, err := parse(ctx, data)

	if err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "the product is found", slog.String("sku", p.Sku))
	}

	return p, err
}

// Search is looking into Redis to find dat matching  the criteria.
// offset are num are coming from Redis api, here is the documentation:
// limits the results to the offset and number of results given.
// Note that the offset is zero-indexed.
// The default is 0 10, which returns 10 items starting from the first result.
// You can use LIMIT 0 0 to count the number of documents in the result set without actually returning them.
func Search(ctx context.Context, q Query, offset, num int) (SearchResults, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "searching products")

	qs := fmt.Sprintf("FT.SEARCH %s \"@status:{online}@type:{product}", db.ProductIdx)

	if q.Keywords != "" {
		k := db.SearchValue(q.Keywords)
		qs += fmt.Sprintf("(@title:%s)|(@description:%s)|(@sku:{%s})|(@id:{%s})", k, k, k, k)
	}

	if q.Slug != "" {
		k := db.SearchValue(q.Slug)
		qs += fmt.Sprintf("(@slug:{%s})", k)
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
		s := SerializeMeta(ctx, q.Meta)
		s = strings.Replace(s, ";", " | ", 999)
		s = db.Escape(s)
		qs += fmt.Sprintf("@meta:{%s}", s)
	}

	qs += fmt.Sprintf("\" SORTBY updated_at desc LIMIT %d %d DIALECT 2", offset, num)

	slog.LogAttrs(ctx, slog.LevelInfo, "preparing redis request", slog.String("query", qs))

	args, err := db.SplitQuery(ctx, qs)
	if err != nil {
		return SearchResults{}, err
	}

	cmds, err := db.Redis.Do(ctx, args...).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot run the search query", slog.String("error", err.Error()))
		return SearchResults{}, err
	}

	res := cmds.(map[interface{}]interface{})
	total := res["total_results"].(int64)

	slog.LogAttrs(ctx, slog.LevelInfo, "search done", slog.Int64("results", total))

	results := res["results"].([]interface{})
	products := []Product{}

	for _, value := range results {
		m := value.(map[interface{}]interface{})
		attributes := m["extra_attributes"].(map[interface{}]interface{})
		data := db.ConvertMap(attributes)

		product, err := parse(ctx, data)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the product", slog.Any("product", data), slog.String("error", err.Error()))
			continue
		}

		products = append(products, product)
	}

	// user, ok := ctx.Value(contexts.User).(users.User)
	// if !ok || user.Role != "admin" {
	// 	tra := map[string]string{
	// 		"query": fmt.Sprintf("'%s'", qs),
	// 	}

	// 	go tracking.Log(ctx, "product_search", tra)
	// }

	return SearchResults{
		Total:    int(total),
		Products: products,
	}, nil
}

func (p Product) URL() string {
	return conf.WebsiteURL + "/" + p.ID + "-" + p.Slug + ".html"
}

// SerializeMeta transforms a meta map to a string representation.
// The values are separated by ";".
// Example: map["color"][]{"blue"} => color_blue
func SerializeMeta(ctx context.Context, m map[string][]string) string {
	slog.LogAttrs(ctx, slog.LevelInfo, "serializing the product meta", slog.Any("meta", m))

	meta := []string{}

	for key, values := range m {
		for _, val := range values {
			meta = append(meta, key+"_"+val)
		}
	}

	s := strings.Join(meta, ";")

	slog.LogAttrs(ctx, slog.LevelInfo, "serialize done successfully", slog.String("serialized", s))

	return s
}

// UnSerializeMeta transform the meta serialized to a map.
// The values are separated by ";".
// Example: color_blue => map["color"]"blue"
func UnSerializeMeta(ctx context.Context, s string) map[string][]string {
	if s == "" {
		return map[string][]string{}
	}

	values := strings.Split(s, ";")
	meta := map[string][]string{}

	for _, value := range values {
		parts := strings.Split(value, "_")

		if len(parts) != 2 {
			slog.LogAttrs(ctx, slog.LevelError, "cannot unserialize the product meta", slog.String("serialized", s))
			continue
		}

		meta[parts[0]] = append(meta[parts[0]], parts[1])
	}

	return meta
}

func Delete(ctx context.Context, pid string) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "deleting product", slog.Any("pid", pid))

	if pid == "" {
		slog.LogAttrs(ctx, slog.LevelInfo, "the pid cannot be empty")
		return errors.New("input:id")
	}

	if _, err := db.Redis.Del(ctx, "product:"+pid).Result(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot delete product", slog.String("string", err.Error()))
		return err
	}

	return nil
}

func List(ctx context.Context, pids []string) ([]Product, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "listing products", slog.Any("pids", pids))

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, key := range pids {
			key := "product:" + key
			rdb.HGetAll(ctx, key)
		}

		return nil
	})

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the products", slog.String("error", err.Error()))
		return []Product{}, err
	}

	pds := []Product{}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the product", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		val := cmd.(*redis.MapStringStringCmd).Val()

		p, err := parse(ctx, val)
		if err != nil {
			continue
		}

		pds = append(pds, p)

	}

	return pds, nil
}
