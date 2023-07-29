// Package parser provides features related to csv parsing
package parser

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/locales"
	"gifthub/products"
	"gifthub/string/stringutil"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/text/message"
)

// lines representes the lines of a csv
type lines [][]string

type csvline struct {
	MID         string
	Sku         string
	Title       string
	Price       float64
	Currency    string
	Quantity    int
	Status      string
	Description string
	Images      []string
	Weight      float64
	Tags        []string
	Links       []string
	Options     map[string]string
}

const isku int = 0
const ititle = 1
const iprice = 2
const icurrency = 3
const iquantity = 4
const istatus = 5
const idescription = 6
const iimages = 7
const iweight = 8
const itags = 9
const ilinks = 10
const ioptions = 11
const cellSeparator = ";"
const optionSeparator = ":"
const requiredFields = 8

var (
	printer = message.NewPrinter(locales.Console)
)

func getUrlExtension(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	position := strings.LastIndex(u.Path, ".")
	if position == -1 {
		log.Printf("input_validation_fail: the url %s does not contain image extension", raw)
		return "", errors.New(printer.Sprintf("csv_image_extension_missing", raw))
	}

	extension := strings.ToLower(u.Path[position+1 : len(u.Path)])

	if !strings.Contains(products.ImageExtensions, extension) {
		log.Printf("input_validation_fail: the extension %s is not supported", extension)
		return "", errors.New(printer.Sprintf("csv_image_extension_not_supported", extension))
	}

	return u.Path[position+1 : len(u.Path)], nil
}

func getFile(url string) (string, error) {
	if strings.HasPrefix(url, "http") {
		return downloadFile(url)
	}

	return copyFile(url)
}

func downloadFile(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("input_validation_fail: the status code %d is not correct", response.StatusCode)
		return "", errors.New(printer.Sprintf("http_bad_status", response.StatusCode, url))
	}

	id, err := stringutil.Random()
	if err != nil {
		log.Printf("ERROR: sequence_fail: something went wrong when generating id %s", err.Error())
		return "", errors.New(printer.Sprintf("something_went_wrong"))
	}

	extension, err := getUrlExtension(url)
	if err != nil {
		return "", err
	}

	p := path.Join(os.TempDir(), fmt.Sprintf("%s.%s", id, extension))
	file, err := os.Create(p)
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return p, nil

}

func copyFile(src string) (string, error) {
	extension := strings.Replace(filepath.Ext(src), ".", "", 1)
	if !strings.Contains(products.ImageExtensions, extension) {
		log.Printf("input_validation_fail: the extension %s is not supported", extension)
		return "", errors.New(printer.Sprintf("csv_image_extension_not_supported", extension))
	}

	stat, err := os.Stat(src)
	if err != nil {
		return "", err
	}

	if !stat.Mode().IsRegular() {
		log.Printf("input_validation_fail: the file %s is not a regular file", src)
		return "", errors.New(printer.Sprintf("csv_bad_file", src))
	}

	s, err := os.Open(src)
	if err != nil {
		return "", err
	}

	defer s.Close()

	id, err := stringutil.Random()
	if err != nil {
		log.Printf("sequence_fail: something went wrong when generating id %s", err.Error())
		return "", errors.New(printer.Sprintf("something_went_wrong"))
	}

	p := path.Join(os.TempDir(), fmt.Sprintf("%s.%s", id, extension))

	d, err := os.Create(p)
	if err != nil {
		return "", err
	}

	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return "", err
	}

	return p, nil
}

func parseCsvLine(line []string) (csvline, error) {
	product := csvline{}

	if len(line) < requiredFields {
		log.Printf("input_validation_fail: csv not valid")
		return product, errors.New(printer.Sprintf("csv_not_valid"))
	}

	sku := line[isku]
	if sku == "" {
		log.Printf("input_validation_fail: the sku is required")
		return product, errors.New(printer.Sprintf("input_required", "sku"))
	}

	isValid := regexp.MustCompile(`^[0-9a-z]+$`).MatchString
	if !isValid(sku) {
		log.Printf("input_validation_fail: the sku %s is invalid", sku)
		return product, errors.New(printer.Sprintf("input_validation", "sku"))
	}

	product.Sku = sku

	title := line[ititle]
	if title == "" {
		log.Printf("input_validation_fail: the title is required")
		return product, errors.New(printer.Sprintf("input_required", "title"))
	}

	product.Title = strings.ReplaceAll(title, "\"", "")

	price, priceErr := strconv.ParseFloat(line[iprice], 32)
	if priceErr != nil {
		log.Printf("input_validation_fail: the price %s is invalid", line[iprice])
		return product, errors.New(printer.Sprintf("input_validation", "price"))
	}

	product.Price = price

	currency := line[icurrency]
	if !conf.IsCurrencySupported(currency) {
		log.Printf("input_validation_fail: the currency %s is invalid", currency)
		return product, errors.New(printer.Sprintf("input_validation", "currency"))
	}

	product.Currency = currency

	quantity, quantityErr := strconv.ParseInt(line[iquantity], 10, 32)
	if quantityErr != nil {
		log.Printf("input_validation_fail: the quantity %s is invalid", line[iquantity])
		return product, errors.New(printer.Sprintf("input_validation", "quantity"))
	}

	product.Quantity = int(quantity)

	status := line[istatus]
	if status != "online" && status != "offline" {
		log.Printf("input_validation_fail: the status %s is invalid", status)
		return product, errors.New(printer.Sprintf("input_validation", "status"))
	}

	product.Status = status

	description := line[idescription]
	if description == "" {
		log.Printf("input_validation_fail: the description is required")
		return product, errors.New(printer.Sprintf("input_required", "description"))
	}

	product.Description = strings.ReplaceAll(description, "\"", "")

	images := strings.Split(line[iimages], ";")
	if len(images) == 0 {
		log.Printf("input_validation_fail: the images are required")
		return product, errors.New(printer.Sprintf("input_required", "images"))
	}

	var paths []string
	for _, v := range images {
		p, err := getFile(v)

		if err != nil {
			log.Printf("input_validation_fail: %s", err.Error())
			return product, errors.New(printer.Sprintf("input_validation", "images"))
		}

		paths = append(paths, p)
	}

	if len(paths) != len(images) {
		log.Printf("input_validation_fail: the images contains error")
		return product, errors.New(printer.Sprintf("input_validation", "images"))
	}

	product.Images = paths

	length := len(line)

	var weight float64
	if length > iweight && line[iweight] != "" {
		w, weightErr := strconv.ParseFloat(line[iweight], 32)

		if weightErr != nil {
			log.Printf("input_validation_fail: the weight %s is not correct", line[iweight])
			return product, errors.New(printer.Sprintf("input_validation", "weight"))
		} else {
			weight = w
		}
	}

	product.Weight = weight

	var tags []string
	if length > itags && line[itags] != "" {
		tags = strings.Split(line[itags], cellSeparator)
	}

	product.Tags = tags

	var links []string
	if length > ilinks && line[ilinks] != "" {
		links = strings.Split(line[ilinks], cellSeparator)
	}

	product.Links = links

	options := make(map[string]string)
	if length > ioptions && line[ioptions] != "" {
		o := strings.Split(line[ioptions], cellSeparator)

		for j, v := range o {
			parts := strings.Split(v, optionSeparator)
			if len(parts) != 2 {
				log.Printf("input_validation_fail: the option %d is not correct %s", j, v)
				return product, errors.New(printer.Sprintf("input_validation", "options"))
			}

			k := strings.ReplaceAll(parts[0], "\"", "")
			v := strings.ReplaceAll(parts[1], "\"", "")

			options[k] = v
		}

		if len(o) != len(options) {
			log.Printf("input_validation_fail: the options contain error")
			return product, errors.New(printer.Sprintf("input_validation", "options"))
		}
	}

	product.Options = options

	return product, nil
}

func deletePreviousImages(ctx context.Context, pid string) error {
	key := "product:" + pid

	v, err := db.Redis.HGet(ctx, key, "images").Result()
	if err != nil {
		log.Printf("ERROR: sequence_fail: redis error %s", err.Error())
		return errors.New(printer.Sprintf("something_went_wrong"))
	}

	img, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("ERROR: sequence_fail: redis error %s", err.Error())
		return errors.New(printer.Sprintf("something_went_wrong"))
	}

	var i int = 0
	for i = 0; i < img; i++ {
		_, p := products.ImagePath(pid, i)

		if err := os.Rename(v, p); err != nil {
			log.Printf("sequence_fail: error when removing %s - %s", p, err.Error())
			return errors.New(printer.Sprintf("something_went_wrong"))
		}
	}

	return nil
}

func createImages(pid string, product csvline) error {
	for i, v := range product.Images {
		folder, p := products.ImagePath(pid, i)

		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			log.Printf("ERROR: sequence_fail: error when moving %s to %s - %s", v, p, err.Error())
			return errors.New(printer.Sprintf("something_went_wrong"))
		}

		if err := os.Rename(v, p); err != nil {
			log.Printf("ERROR: sequence_fail: error when moving %s to %s-  %s", v, p, err.Error())
			return errors.New(printer.Sprintf("something_went_wrong"))
		}
	}

	return nil
}

func removeTmpFiles(product csvline) {
	if len(product.Images) > 0 {
		for _, v := range product.Images {
			if err := os.Remove(v); err != nil {
				log.Printf("WARN: sequence_fail: error when removing %s file %s", v, err.Error())
			}
		}
	}
}

func addProductToRedis(ctx context.Context, pid string, product csvline) error {
	key := "product:" + pid
	lkey := "product:links:" + pid
	tkey := "product:tags:" + pid
	okey := "product:options:" + pid
	now := time.Now()

	_, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, lkey, okey, tkey)
		rdb.HSet(ctx, key,
			"pid", pid,
			"sku", product.Sku,
			"title", product.Title,
			"description", product.Description,
			"images", len(product.Images),
			"currency", product.Currency,
			"price", product.Price,
			"quantity", product.Quantity,
			"status", product.Status,
			"weight", product.Weight,
			"mid", product.MID,
			"updated_at", time.Now().Format(time.RFC3339),
		)
		rdb.ZAdd(ctx, "products:"+product.MID, redis.Z{
			Score:  float64(now.Unix()),
			Member: pid,
		})

		if len(product.Links) > 0 {
			rdb.SAdd(ctx, lkey, product.Links)
		}

		if len(product.Tags) > 0 {
			rdb.SAdd(ctx, tkey, product.Tags)
		}

		if len(product.Options) > 0 {
			for k, v := range product.Options {
				rdb.HSet(ctx, okey, k, v)
			}
		}

		return nil
	})

	return err
}

func processLine(chans chan<- int, i int, mid string, line []string) {
	product, err := parseCsvLine(line)
	if err != nil {
		log.Printf(printer.Sprintf("csv_line_error", i, err.Error()))
		chans <- 0
		return
	}

	catchError := func(err error) {
		log.Printf(printer.Sprintf("csv_line_error", i, err.Error()))
		removeTmpFiles(product)
		chans <- 0
	}

	product.MID = mid

	ctx := context.Background()
	key := "merchant:" + mid + ":" + product.Sku

	pid, err := db.Redis.Get(ctx, key).Result()
	if pid == "" || err != nil {
		pid, err = stringutil.Random()
		if err != nil {
			catchError(err)
			return
		}

		if _, err := db.Redis.Set(ctx, key, pid, 0).Result(); err != nil {
			catchError(err)
			return
		}
	} else {
		deletePreviousImages(ctx, pid)
	}

	if err = createImages(pid, product); err != nil {
		catchError(err)
		return
	}

	if err = addProductToRedis(ctx, pid, product); err != nil {
		catchError(err)
		return
	}

	chans <- 1
}

// Import imports csv data into redis.
// The first line has to contain the headers ordered as the following list:
//   - sku: the unique reference (per merchant)
//   - title: the product title
//   - price: the product price
//   - currency: the product currency
//   - quantity: the product quantity
//   - status: the product status "online" or "offline"
//   - description: the product description
//   - images: the product images
//   - weight: the product weight in grams (optional)
//   - tags: the product tags or categories (optional)
//   - links: the product ids linked to the product
//   - options: the product options. An option is a couple name/value separated by ":"
//
// The separator used inside a cell is ";". If a cell contains a comma, if has to be surrounded
// by double quotes.
// The import will ignore lines if:
//   - the line does have the minimal required values
//   - the line has bad data, for example, if status contains "toto"
//   - one of the images cannot be imported
//   - the currency is not recognized
//   - the field number parsing is invalid
//
// If a product link references a non existing product id, it will be ignored when the
// product details will be displayed.
func Import(data [][]string, mid string) (int, error) {
	// ctx := context.Background()
	lines := 0
	chans := make(chan int, len(data)-1)

	for i, line := range data {
		if i == 0 {
			if len(line) < requiredFields {
				log.Printf("input_validation_fail: csv not valid")
				return 0, errors.New(printer.Sprintf("csv_not_valid"))
			}

			if line[0] != "sku" {
				log.Printf("input_validation_fail: csv header not valid")
				return 0, errors.New(printer.Sprintf("csv_not_valid"))
			}

			continue
		}

		go processLine(chans, i, mid, line)
	}

	for i := 0; i < cap(chans); i++ {
		lines += <-chans
	}

	return lines, nil
}
