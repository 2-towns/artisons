// Package parser provides features related to csv parsing
package parser

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"gifthub/locales"
	"gifthub/products"
	"gifthub/string/stringutil"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// lines representes the lines of a csv
type lines [][]string

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
	slog.Info("get file extension", slog.String("url", raw))

	u, err := url.Parse(raw)
	if err != nil {
		slog.Error("cannot parse the url", slog.String("error", err.Error()))
		return "", err
	}

	position := strings.LastIndex(u.Path, ".")
	if position == -1 {
		slog.Info("cannot proceed the image without extension", slog.String("url", raw))
		return "", errors.New(printer.Sprintf("the image extension is missing in %s", raw))
	}

	extension := strings.ToLower(u.Path[position+1 : len(u.Path)])

	if !strings.Contains(products.ImageExtensions, extension) {
		slog.Info("cannot proceed the unsupported image extension", slog.String("url", raw))
		return "", errors.New(printer.Sprintf("the image extension %s is not supported", extension))
	}

	slog.Info("file extension detected", slog.String("extension", extension))

	return extension, nil
}

func getFile(url string) (string, error) {
	if strings.HasPrefix(url, "http") {
		return downloadFile(url)
	}

	return copyFile(url)
}

func downloadFile(url string) (string, error) {
	slog.Info("downloading the file", slog.String("url", url))

	response, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		slog.Info("cannot download the file", slog.String("url", url), slog.Int("status_code", response.StatusCode))
		return "", errors.New(printer.Sprintf("error_http_badstatus", response.StatusCode, url))
	}

	id, err := stringutil.Random()
	if err != nil {
		slog.Error("cannot generated the id", slog.String("error", err.Error()))
		return "", errors.New(printer.Sprintf("something went wrong"))
	}

	extension, err := getUrlExtension(url)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.%s", id, extension)
	p := path.Join(conf.WorkingSpace, "web", "tmp", filename)
	file, err := os.Create(p)
	if err != nil {
		slog.Error("cannot generated local folder", slog.String("error", err.Error()), slog.String("extension", extension), slog.String("id", id))
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		slog.Error("cannot copy the file", slog.String("error", err.Error()))
		return "", err
	}

	slog.Info("file downloaded", slog.String("url", url), slog.String("id", id))

	return filename, nil

}

func copyFile(src string) (string, error) {
	slog.Info("copying the file", slog.String("src", src))

	extension := strings.Replace(filepath.Ext(src), ".", "", 1)
	if !strings.Contains(products.ImageExtensions, extension) {
		slog.Info("cannot support the extension", slog.String("extension", extension))
		return "", errors.New(printer.Sprintf("the image extension %s is not supported", extension))
	}

	stat, err := os.Stat(src)
	if err != nil {
		return "", err
	}

	if !stat.Mode().IsRegular() {
		slog.Info("cannot copy a non regular file", slog.String("src", src))
		return "", errors.New(printer.Sprintf("error_csv_badfile", src))
	}

	s, err := os.Open(src)
	if err != nil {
		slog.Error("cannot copy the file", slog.String("error", err.Error()), slog.String("src", src))
		return "", err
	}

	defer s.Close()

	id, err := stringutil.Random()
	if err != nil {
		slog.Error("cannot generated the id", slog.String("error", err.Error()), slog.String("src", src))
		return "", errors.New(printer.Sprintf("something went wrong"))
	}

	filename := fmt.Sprintf("%s.%s", id, extension)
	p := path.Join(conf.WorkingSpace, "web", "tmp", filename)

	d, err := os.Create(p)
	if err != nil {
		slog.Error("cannot create the path", slog.String("path", p), slog.String("error", err.Error()), slog.String("src", src))
		return "", err
	}

	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		slog.Error("cannot copy the path", slog.String("path", p), slog.String("src", src), slog.String("error", err.Error()))
		return "", err
	}

	slog.Info("file copied successfully", slog.String("path", p), slog.String("id", id))

	return filename, nil
}

func parseCsvLine(line []string) (products.Product, error) {
	slog.Info("parsing the line", slog.String("src", strings.Join(line, ",")))

	if len(line) < requiredFields {
		slog.Info("cannot parse the csv", slog.Int("length", len(line)), slog.Int("required", requiredFields))
		return products.Product{}, errors.New(printer.Sprintf("the csv is invalid"))
	}

	price, priceErr := strconv.ParseFloat(line[iprice], 32)
	if priceErr != nil {
		slog.Error("cannot parse the price", slog.Float64("price", price), slog.String("error", priceErr.Error()))
		return products.Product{}, errors.New(printer.Sprintf("input:price", "price"))
	}

	quantity, quantityErr := strconv.ParseInt(line[iquantity], 10, 32)
	if quantityErr != nil {
		slog.Error("cannot parse the quantity", slog.Int64("quantity", quantity), slog.String("error", quantityErr.Error()))
		return products.Product{}, errors.New(printer.Sprintf("input:quantity", "quantity"))
	}

	images := strings.Split(line[iimages], ";")
	if len(images) == 0 {
		slog.Info("cannot parse the empty images")
		return products.Product{}, errors.New("input:image_1")
	}

	var iids []string
	for index, v := range images {
		iid, err := getFile(v)
		if err != nil {
			slog.Error("cannot get the file", slog.String("error", err.Error()))
			return products.Product{}, errors.New(printer.Sprintf("input:image_%d", index))
		}

		iids = append(iids, iid)
	}

	if len(iids) != len(images) {
		slog.Info("cannot parse the images", slog.Int("length", len(iids)), slog.Int("images", len(images)))
		return products.Product{}, errors.New(printer.Sprintf("input:images", "images"))
	}

	length := len(line)

	var weight float64
	if length > iweight && line[iweight] != "" {
		w, weightErr := strconv.ParseFloat(line[iweight], 32)

		if weightErr != nil {
			slog.Error("cannot parse the weight", slog.String("weight", line[iweight]), slog.String("error", weightErr.Error()))
			return products.Product{}, errors.New(printer.Sprintf("input:weight", "weight"))
		} else {
			weight = w
		}
	}

	var tags []string
	if length > itags && line[itags] != "" {
		tags = strings.Split(line[itags], cellSeparator)
	}

	var links []string
	if length > ilinks && line[ilinks] != "" {
		links = strings.Split(line[ilinks], cellSeparator)
	}

	options := make(map[string]string)
	if length > ioptions && line[ioptions] != "" {
		o := strings.Split(line[ioptions], cellSeparator)

		for j, v := range o {
			parts := strings.Split(v, optionSeparator)
			if len(parts) != 2 {
				slog.Info("cannot parse the option", slog.Int("index", j), slog.String("option", v))
				return products.Product{}, errors.New(printer.Sprintf("input:options", "options"))
			}

			k := strings.ReplaceAll(parts[0], "\"", "")
			v := strings.ReplaceAll(parts[1], "\"", "")

			options[k] = v
		}

		if len(o) != len(options) {
			slog.Info("cannot parse the options", slog.Int("length", len(o)), slog.Int("options", len(options)))
			return products.Product{}, errors.New(printer.Sprintf("input:options", "options"))
		}
	}

	product := products.Product{
		Sku:         line[isku],
		Title:       strings.ReplaceAll(line[ititle], "\"", ""),
		Description: strings.ReplaceAll(line[idescription], "\"", ""),
		Price:       price,
		Currency:    line[icurrency],
		Quantity:    int(quantity),
		Status:      line[istatus],
		Weight:      weight,
		Tags:        tags,
		Links:       links,
		Meta:        options,
		Image1:      iids[0],
	}

	ctx := context.WithValue(context.Background(), contexts.Locale, language.English)
	if err := product.Validate(ctx); err != nil {
		return products.Product{}, err
	}

	slog.Info("line parsed successfully")

	return product, nil
}

func deletePreviousImages(ctx context.Context, pid string) error {
	l := slog.With(slog.String("id", pid))
	l.Info("deleting previous images")

	key := "product:" + pid

	v, err := db.Redis.HGet(ctx, key, "images").Result()
	if err != nil {
		l.Error("cannot store the product", slog.String("error", err.Error()))
		return errors.New(printer.Sprintf("something went wrong"))
	}

	img, err := strconv.Atoi(v)
	if err != nil {
		l.Error("cannot convert the index", slog.String("index", v), slog.String("error", err.Error()))
		return errors.New(printer.Sprintf("something went wrong"))
	}

	for i := 0; i < img; i++ {
		p := products.ImagePath(pid)

		if err := os.Rename(v, p); err != nil {
			l.Error("cannot remove the image", slog.String("path", p), slog.String("error", err.Error()))
			return errors.New(printer.Sprintf("something went wrong"))
		}
	}

	l.Info("previous images deleted")

	return nil
}

func createImages(product products.Product) error {
	l := slog.With(slog.String("id", product.ID))
	l.Info("creating previous images")

	if product.Image1 != "" {
		p := products.ImagePath(product.Image1)
		old := path.Join(conf.WorkingSpace, "web", "tmp", product.Image1)
		if err := os.Rename(old, p); err != nil {
			l.Error("cannot move the file", slog.String("old", old), slog.String("new", p), slog.String("error", err.Error()))
			return errors.New(printer.Sprintf("something went wrong"))
		}
	}

	if product.Image2 != "" {
		p := products.ImagePath(product.Image2)
		old := path.Join(conf.WorkingSpace, "web", "tmp", product.Image2)
		if err := os.Rename(old, p); err != nil {
			l.Error("cannot move the file", slog.String("old", old), slog.String("new", p), slog.String("error", err.Error()))
			return errors.New(printer.Sprintf("something went wrong"))
		}
	}

	if product.Image3 != "" {
		p := products.ImagePath(product.Image3)
		old := path.Join(conf.WorkingSpace, "web", "tmp", product.Image3)
		if err := os.Rename(old, p); err != nil {
			l.Error("cannot move the file", slog.String("old", old), slog.String("new", p), slog.String("error", err.Error()))
			return errors.New(printer.Sprintf("something went wrong"))
		}
	}
	if product.Image4 != "" {
		p := products.ImagePath(product.Image4)
		old := path.Join(conf.WorkingSpace, "web", "tmp", product.Image4)
		if err := os.Rename(old, p); err != nil {
			l.Error("cannot move the file", slog.String("old", old), slog.String("new", p), slog.String("error", err.Error()))
			return errors.New(printer.Sprintf("something went wrong"))
		}
	}

	l.Info("previous images created")

	return nil
}

func removeTmpFiles(product products.Product) {
	slog.Info("remove temporary images")

	if product.Image1 != "" {
		if err := os.Remove(product.Image1); err != nil {
			slog.Error("cannot move the file", slog.String("old", product.Image1), slog.String("error", err.Error()))
		}
	}

	if product.Image2 != "" {
		if err := os.Remove(product.Image2); err != nil {
			slog.Error("cannot move the file", slog.String("old", product.Image2), slog.String("error", err.Error()))
		}
	}

	if product.Image3 != "" {
		if err := os.Remove(product.Image3); err != nil {
			slog.Error("cannot move the file", slog.String("old", product.Image3), slog.String("error", err.Error()))
		}
	}

	if product.Image4 != "" {
		if err := os.Remove(product.Image4); err != nil {
			slog.Error("cannot move the file", slog.String("old", product.Image4), slog.String("error", err.Error()))
		}
	}

	slog.Info("temporary images deleted")
}

func processLine(chans chan<- int, i int, mid string, line []string) {
	l := slog.With(slog.Int("index", i))
	l.Info("processing the file", slog.String("mid", mid))

	product, err := parseCsvLine(line)
	if err != nil {
		l.Error("cannot parse the csv line", slog.String("error", err.Error()))
		chans <- 0
		return
	}

	catchError := func(err error) {
		l.Error("cannot parse the csv", slog.String("error", err.Error()))
		removeTmpFiles(product)
		chans <- 0
	}

	product.MID = mid

	ctx := context.WithValue(context.Background(), contexts.Locale, language.English)
	key := "merchant:" + mid + ":" + product.Sku

	exists, err := db.Redis.Exists(ctx, key).Result()
	if err != nil {
		catchError(err)
		return
	}

	var pid string

	if exists == 0 {
		pid, err = stringutil.Random()
		if err != nil {
			catchError(err)
			return
		}
	} else {
		pid, err = db.Redis.Get(ctx, key).Result()
		if err != nil {
			catchError(err)
			return
		}
	}

	if pid != "" {
		deletePreviousImages(ctx, pid)
	} else {
		pid, err = stringutil.Random()
		if err != nil {
			catchError(err)
			return
		}
	}

	product.ID = pid

	if err = createImages(product); err != nil {
		catchError(err)
		return
	}

	if err = product.Save(ctx); err != nil {
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
	slog.Info("importing the data", slog.String("mid", mid))

	lines := 0
	chans := make(chan int, len(data)-1)

	for i, line := range data {
		if i == 0 {
			if len(line) < requiredFields {
				slog.Info("cannot parse the csv line", slog.Int("index", i), slog.Int("len", len(line)), slog.Int("required", requiredFields))
				return 0, errors.New(printer.Sprintf("the csv is invalid"))
			}

			if line[0] != "sku" {
				slog.Info("cannot parse the sku empty", slog.Int("index", i))
				return 0, errors.New(printer.Sprintf("the csv is invalid"))
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
