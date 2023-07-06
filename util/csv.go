package util

import (
	"errors"
	"log"
)

const isku = 0
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

const requiredFields = 8

// CsvImport imports csv data into redis.
// The first line has to contain the headers ordered as the following list:
//   - sku: the unique reference (per merchant)
//   - title: the product title
//   - price: the product price
//   - currency: the product currency
//   - quantity: the product quantity
//   - status: the product status "online" or "offline"
//   - description: the product description
//   - images: the product images
//   - weight: the product weight (optional)
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
func CsvImport(data CsvLines) (int, error) {
	// ctx := context.Background()
	lines := 0

	for i, line := range data {
		if i == 0 {
			if len(line) < requiredFields {
				log.Println("input_validation_fail", "The CSV is not valid")

				return 0, errors.New("input_validation_fail: csv not valid")
			}

			if line[0] != "sku" {
				log.Println("input_validation_fail", "The CSV header is not valid")

				return 0, errors.New("input_validation_fail: csv header not valid")
			}
		} else {
			// price, err := strconv.ParseFloat(line[indexes["price"]], 32)

			// if err != nil {
			// 	log.Println(err)

			// 	continue
			// }

			// id := line[indexes["id"]]

			// isValid := regexp.MustCompile(`^[0-9a-z]+$`).MatchString
			// if !isValid(id) {
			// 	log.Printf("The id %s is not valid.", id)

			// 	continue
			// }

			// slug := Slugify(id + " " + line[indexes["title"]])

			// // todo validation
			// log.Println(indexes)

			// key := "product:" + id
			// if _, err := RedisClient.Pipelined(ctx, func(rdb redis.Pipeliner) error {
			// 	rdb.HSet(ctx, key, "id", id)
			// 	rdb.HSet(ctx, key, "title", line[indexes["title"]])
			// 	rdb.HSet(ctx, key, "image", line[indexes["image"]])
			// 	rdb.HSet(ctx, key, "description", line[indexes["description"]])
			// 	rdb.HSet(ctx, key, "price", price)
			// 	rdb.HSet(ctx, key, "slug", slug)

			// 	options := 0
			// 	i := oIndex
			// 	length := len(line)
			// 	k := key + ":options"

			// 	for i+1 < length {
			// 		if (line[i] != "") && (line[i+1] != "") {
			// 			rdb.HSet(ctx, k, fmt.Sprintf("option_name_%d", options), line[i])
			// 			rdb.HSet(ctx, k, fmt.Sprintf("option_value_%d", options), line[i+1])
			// 			options++
			// 		}

			// 		i++
			// 	}

			// 	if options > 0 {
			// 		rdb.HSet(ctx, k, "options", options)
			// 	}

			// 	return nil
			// }); err != nil {
			// 	panic(err)
			// }*

			lines++
		}
	}

	return lines, nil
}
