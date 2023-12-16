// Package products provide everything around products
package products

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/tracking"
	"gifthub/users"
	"gifthub/utils"
	"log/slog"
	"math/rand"
	"time"

	"gifthub/db"
	"log"
	"strconv"
	"strings"

	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator/v10"

	"github.com/redis/go-redis/v9"
)


type Product struct {

	ID          string  `redis:"id"` // ID is an unique identifier
	Title       string  `redis:"title" validate:"required"`
	Description string  `redis:"description" validate:"required"`
	Price       float64 `redis:"price"`
	// The percent discount
	Discount float64 `redis:"discount"`
	Slug     string  `redis:"slug"`
	MID      string  `redis:"mid"`
	Sku      string  `redis:"sku" validate:"required,alphanum"`
	Currency string  `redis:"currency"`
	Quantity int     `redis:"quantity"`
	// Images length
	Length int     `redis:"length" validate:"required"`
	Status string  `redis:"status" validate:"oneof=online offline"`
	Weight float64 `redis:"weight"`

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
	ID          string            `redis:"id"`
	MID         string            `redis:"merchant_id"` // MerchantID is the id of the merchant that sells the product
	Title       string            `redis:"title"`
	Image       string            `redis:"image"`
	Description string            `redis:"description"`
	Price       float64           `redis:"price"`
	Slug        string            `redis:"slug"`
	Links       []string          // Links contains the linked product IDs
	Meta        map[string]string // Meta contains the product options.

}

// Utilise le type MetaData pour les options.


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

//Add new product to the database
func Add(product Product) error {
	v := validator.New()

	// Validate the product
	if err := v.Struct(product); err != nil {
    slog.Error("input_validation_fail", "error", err.Error(),"product", product)
    return errors.New("product_invalid")
}

	ctx := context.Background()

	// Generating a new unique ID for the product
	pid, err := utils.RandomString(10)
	if err != nil {
		slog.Error("sequence_fail", "error", err.Error(), "description", "got error from Redis while generating random string")
		return errors.New("something_went_wrong")
	}


	// Adding the ID to the product structure
	product.ID = pid

	// Add product to Redis
	pipe := db.Redis.Pipeline()
	score, err := db.Redis.Incr(ctx, "product:score").Result()

	// Update products list
	pipe.ZAdd(ctx, "products", redis.Z{Score: float64(score), Member: pid})

	// Store product data
	pipe.HSet(ctx, fmt.Sprintf("product:%s", pid), map[string]interface{}{
			"id":          product.ID,
			"title":       product.Title,
			"image":       product.Images,
			"description": product.Description,
			"price":       strconv.FormatFloat(product.Price, 'f', -1, 64),
			"slug":        product.Slug,
	})

pipe.HSet(ctx, fmt.Sprintf("product:%s", pid), product)

		if _,err = pipe.Exec(ctx); err!=nil {
		log.Printf("ERROR: sequence_fail: go error from redis %s", err.Error())
		return errors.New("something_went_wrong")
	}

	log.Printf("WARN: sensitive_create: a new product is created with id %s\n", pid)

	return nil
}

// func parseProduct(product Product) (Product, error) {
// 	if product.ID == "" {
// 		return Product{}, errors.New("ID is missing")
// 	}

// 	var merchantID string
// 	if m["merchant_id"] != "" {
// 		merchantID = m["merchant_id"]
// 	}

// 	if product.Price == 0 {
//     slog.Error("sequence_fail", "error", "Price is missing or zero", "product", product)
//     return Product{}, errors.New("price_is_missing_or_zero")
// }

// price := math.Round(product.Price * 100) / 100


// 	// Valider ou transformer l'ID du marchand
// 	if product.MID == "" {
// 		return Product{}, errors.New("merchant_id is missing")
// 	}


// 	// Mettre en majuscule la première lettre de chaque mot du titre
// 	t := cases.Title(language.English)
// 	// Mettre en majuscule la première lettre du titre
// 	product.Title = t.String(product.Title)

// 	// Retirer les espaces inutiles dans la description
// 	product.Description = strings.TrimSpace(product.Description)


// 	return Product{

// 		ID:          strconv.FormatInt(id, 10),
// 		Title:       product.Title,
// 		Image:       product.Image,
// 		Description: product.Description,
// 		Price:       price,
// 		Slug:        product.Slug,
// 		MID:         product.MID,
// 		Links:       product.Links,
// 		Meta:        product.Meta,
// 	}, nil
// }




func List(page int64) ([]Product, error) {
	key := "products"
	ctx := context.Background()

	var start, end int64
	if page == -1 {
			start = 0
			end = -1
	} else {
			start = page * conf.ItemsPerPage
			end = start + conf.ItemsPerPage - 1
	}

	products := []Product{}
	ids := db.Redis.ZRange(ctx, key, start, end).Val()
	pipe := db.Redis.Pipeline()

	for _, id := range ids {
			k := "product:" + id
			pipe.HGetAll(ctx, k)
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
			log.Printf("ERROR: sequence_fail: go error from redis %s", err.Error())
			return nil, err
	}

	for _, cmd := range cmds {
			cmdResult := cmd.(*redis.Cmd)
			m, err := cmdResult.Result()
			if err != nil {
					slog.Error("command_result_error", "error", err.Error(), "description", "Error getting command result")
					continue
			}

			productMap, ok := m.(map[string]interface{})
			if !ok {
					log.Println("Type assertion to map[string]interface{} failed")
					continue
			}

			stringMap := make(map[string]string)
			for k, v := range productMap {
					val, ok := v.(string)
					if ok {
							stringMap[k] = val
					}
			}

			product, err := parse(ctx, stringMap)
			if err != nil {
					log.Printf("ERROR: failed to parse product: %s", err)
					continue
			}

			products = append(products, product)
	}

	return products, nil
}




	func Fake() Product {
		randID, err := utils.RandomString(10)
		if err != nil {
			log.Printf("ERROR: sequence_fail: error when generating random ID %s", err.Error())
			return Product{}
		}
		randPrice := rand.Float64() * 100
		return Product{
			ID:          randID,                                                            // Génère un ID aléatoire entre 1 et 1000.
			Title:       faker.Sentence(),                                                  // Génère une phrase aléatoire pour le titre.
			Description: faker.Paragraph(),                                                 // Génère un paragraphe aléatoire pour la description.
			Price:       float64(randPrice),                                                         // Génère un prix aléatoire.
			Slug:        faker.Word(),                                                      // Génère un mot aléatoire pour le slug.
			Images:       []string{"https://example.com/image.jpg","https://example.com/image.jpg"},                             // Utilise une URL d'image statique.
			MID:         faker.UUIDHyphenated(),                                            // Génère une chaîne aléatoire.
			Links:       []string{faker.URL(), faker.URL()},                                // Génère deux URLs aléatoires pour les liens.
			Meta:        map[string]string{"key": faker.Word(), "value": faker.Sentence()}, // Génère un dictionnaire avec une clé et une valeur aléatoires.
		}
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
	slog.LogAttrs(c, slog.LevelInfo, "parsing the product data")

	l, err := strconv.ParseInt(data["length"], 10, 32)
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

	quantity, err := strconv.ParseInt(data["quantity"], 10, 32)
	if err != nil {
		slog.Error("cannot parse the product quantity", slog.String("quantity", data["quantity"]))
		return Product{}, err
	}

	var weight float64

	if data["weight"] != "" {
		v, err := strconv.ParseFloat(data["weight"], 32)
		if err != nil {
			slog.Error("cannot parse the product weight", slog.String("weight", data["weight"]))
			return Product{}, err
		}

		weight = v
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
		Price:       price,
		Slug:        data["slug"],
		MID:         data["mid"],
		Sku:         data["sku"],
		Currency:    data["currency"],
		Quantity:    int(quantity),
		Weight:      weight,
		Status:      data["status"],
		Tags:        strings.Split(data["tags"], ";"),
		Links:       strings.Split(data["links"], ";"),
		Meta:        UnSerializeMeta(c, data["meta"], ";"),
		Length:      length,
		Images:      images,
	}, nil
}

func (p Product) Validate(c context.Context) error {
	slog.LogAttrs(c, slog.LevelInfo, "validating a product")

	v := validator.New()
	if err := v.Struct(p); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot validate the user", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("product_%s_invalid", low)
	}

	if !conf.IsCurrencySupported(p.Currency) {
		slog.Info("cannot use an unsupported currency", slog.String("currency", p.Currency))
		return errors.New("product_currency_invalid")
	}

	return nil
}

// Save a product into redis.
// The keys are :
// product:pid => the product data
func (p Product) Save(ctx context.Context) error {
	if p.ID == "" {
		slog.Error("cannot continue with empty pid")
		return errors.New("input_pid_required")
	}

	l := slog.With(slog.String("id", p.ID))
	l.Info("storing the product")

	key := "product:" + p.ID
	// now := time.Now()

	_, err := db.Redis.HSet(ctx, key,
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
		"created_at", time.Now().Unix(),
		"updated_at", time.Now().Unix(),
	).Result()

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

func Search(c context.Context, q Query) ([]Product, error) {
	slog.LogAttrs(c, slog.LevelInfo, "searching products")

	qs := "@status:{online} "

	if q.Keywords != "" {
		qs += fmt.Sprintf("(@title:*%s*)|(@description:*%s*)|(@sku:{%s})", q.Keywords, q.Keywords, q.Keywords)
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
		db.ProductIdx,
		qs,
	).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot run the search query", slog.String("query", qs), slog.String("error", err.Error()))
		return []Product{}, err
	}

	res := cmds.(map[interface{}]interface{})
	slog.LogAttrs(c, slog.LevelInfo, "search done", slog.Int64("results", res["total_results"].(int64)))

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

	user, ok := ctx.Value(contexts.User).(users.User)
	if !ok || user.Role != "admin" {
		tra := map[string]string{
			"query": fmt.Sprintf("'%s'", qs),
		}

		go tracking.Log(c, "product_search", tra)
	}

	return products, nil
}

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

