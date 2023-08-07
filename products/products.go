// Package products provide everything around products
package products

import (
	"context"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"log"

	"github.com/redis/go-redis/v9"
)

// Product is the product representation in the application
type Product struct {
	ID          string            `redis:"id"` // ID is an unique identifier
	Title       string            `redis:"title"`
	Image       string            `redis:"image"`
	Description string            `redis:"description"`
	Price       float64           `redis:"price"`
	Slug        string            `redis:"slug"`
	Links       []string          // Links contains the linked product IDs
	Meta        map[string]string // Meta contains the product options.
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
func Availables(pids []string) bool {
	ctx := context.Background()
	pipe := db.Redis.Pipeline()
	for _, pid := range pids {
		pipe.HGet(ctx, "product:"+pid, "status")
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("sequence_fail: error when getting product ids %v, %v", pids, err.Error())
		return false
	}

	for _, cmd := range cmds {
		status := cmd.(*redis.StringCmd).Val()
		if status != "online" {
			log.Printf("input_validation_fail: error the product %s is not available, %v", cmd.Args()[1], err.Error())
			return false
		}
	}

	return true
}

// Available return true if the product is available
func Available(pid string) bool {
	ctx := context.Background()
	status, err := db.Redis.HGet(ctx, "product:"+pid, "status").Result()
	if err != nil {
		log.Printf("input_validation_fail: error when checking product:%s exists %s", pid, err.Error())
		return false
	}

	return status == "online"
}
