package admin

import (
	"context"
	"errors"
	"gifthub/conf"
	"gifthub/http/httpext"
	"gifthub/products"
	"gifthub/string/stringutil"
	"log/slog"
	"mime/multipart"
	"strconv"
	"strings"
)

func processProductFrom(ctx context.Context, form multipart.Form, id string) (products.Product, error) {
	if len(form.Value["price"]) == 0 || form.Value["price"][0] == "" {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the empty price")
		return products.Product{}, errors.New("input:price")
	}

	if len(form.Value["quantity"]) == 0 || form.Value["quantity"][0] == "" {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the empty quantity")
		return products.Product{}, errors.New("input:quantity")
	}

	price, err := strconv.ParseFloat(form.Value["price"][0], 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the price", slog.String("price", form.Value["price"][0]), slog.String("error", err.Error()))
		return products.Product{}, errors.New("input:price")
	}

	quantity, err := strconv.ParseInt(form.Value["quantity"][0], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the quantity", slog.String("quantity", form.Value["quantity"][0]), slog.String("error", err.Error()))
		return products.Product{}, errors.New("input:quantity")
	}

	var discount float64 = 0
	if len(form.Value["discount"]) > 0 && form.Value["discount"][0] != "" {
		val, err := strconv.ParseFloat(form.Value["discount"][0], 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the discount", slog.String("discount", form.Value["discount"][0]), slog.String("error", err.Error()))
			return products.Product{}, errors.New("input:discount")
		}
		discount = val
	}

	var weight float64 = 0
	if len(form.Value["weight"]) > 0 && form.Value["weight"][0] != "" {
		val, err := strconv.ParseFloat(form.Value["weight"][0], 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the weight", slog.String("weight", form.Value["weight"][0]), slog.String("error", err.Error()))
			return products.Product{}, errors.New("input:weight")
		}
		weight = val
	}

	tags := []string{}
	if len(form.Value["tags"]) > 0 && form.Value["tags"][0] != "" {
		tags = strings.Split(form.Value["tags"][0], " ")
	}

	exists := id != ""
	pid := id
	if pid == "" {
		pid, err = stringutil.Random()
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot generated the product id", slog.String("error", err.Error()))
			return products.Product{}, errors.New("something went wrong")
		}
	}

	sku := ""
	if len(form.Value["sku"]) > 0 {
		sku = form.Value["sku"][0]
	}

	title := ""
	if len(form.Value["title"]) > 0 {
		title = form.Value["title"][0]
	}

	description := ""
	if len(form.Value["description"]) > 0 {
		description = form.Value["description"][0]
	}

	status := ""
	if len(form.Value["status"]) > 0 {
		status = form.Value["status"][0]
	}

	p := products.Product{
		ID:          pid,
		Title:       title,
		Description: description,
		Sku:         sku,
		Status:      status,
		Price:       price,
		Discount:    discount,
		Quantity:    int(quantity),
		Weight:      weight,
		Tags:        tags,
		Currency:    conf.Currency,
	}

	err = p.Validate(ctx)
	if err != nil {
		return products.Product{}, err
	}

	files, err := httpext.ProcessFiles(ctx, form.File, []string{"image_1", "image_2", "image_3", "image_4"})
	if err != nil {
		return products.Product{}, err
	}

	if files["image_1"] == nil && !exists {
		slog.LogAttrs(ctx, slog.LevelInfo, "at least one image is required")
		return products.Product{}, errors.New("input:image_1")
	}

	images, err := httpext.Upload(ctx, files)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot update the files", slog.String("error", err.Error()))
		return products.Product{}, errors.New("something went wrong")
	}

	if images["image_1"] != "" {
		p.Image1 = images["image_1"]
	}

	del2 := form.Value["image_2_delete"]
	if images["image_2"] != "" {
		p.Image2 = images["image_2"]
	} else if len(del2) > 0 && del2[0] != "" {
		p.Image2 = "-"
	}

	del3 := form.Value["image_3_delete"]
	if images["image_3"] != "" {
		p.Image3 = images["image_3"]
	} else if len(del3) > 0 && del3[0] != "" {
		p.Image3 = "-"
	}

	del4 := form.Value["image_4_delete"]
	if images["image_4"] != "" {
		p.Image4 = images["image_4"]
	} else if len(del4) > 0 && del4[0] != "" {
		p.Image4 = "-"
	}

	return p, nil

}
