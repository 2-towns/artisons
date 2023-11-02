// Package products provide everything around products
package products

import (
	"context"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"log/slog"

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
	status, err := db.Redis.HGet(ctx, "product:"+pid, "status").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "the product does not exist", slog.String("error", err.Error()))
		return false
	}

	l.LogAttrs(c, slog.LevelInfo, "got the product status", slog.String("availability", "status"))

	return status == "online"
}
