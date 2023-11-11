package populate

import (
	"context"
	"gifthub/db"
	"gifthub/shops"
	"time"
)

func Shop(ctx context.Context) (shops.Shop, error) {
	now := time.Now()

	_, err := db.Redis.HSet(ctx, "shop",
		"logo", "../web/images/logo",
		"slug", "manger-de-l-ail-c-est-bon-pour-la-santé",
		"address_firstname", "Arnaud",
		"address_lastname", "None",
		"address_city", "Oran",
		"address_street", "Hay Yasmine",
		"address_complementary", "Hay Salam",
		"address_zipcode", "31244",
		"address_phone", "0559682532",
		"updated_at", now.Unix(),
	).Result()

	if err != nil {
		return shops.Shop{}, err
	}

	return shops.Shop{}, err
}
