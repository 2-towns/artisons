// tests gather test utilites
package tests

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/string/stringutil"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/text/language"
)

const ArticleID = 1
const ArticleTitle = "Mangez de l'ail !"
const ArticleDescription = "C'est un antiseptique."
const ArticleCmsID = 2
const CartID = "CAR1"
const CartProductID = "PDT3"
const ProductID1 = "PDT1"
const ProductSKU = "SKU1"
const ProductID2 = "PDT2"
const ProductSlug = "t-shirt-tester-c-est-douter"
const ProductTag = "clothes"
const UserID1 = 1
const UserToDeleteID = 2
const Otp = "123456"
const UserSIDSignedIn = "987654321"
const UserSID = "123456789"
const UserRefreshSID = "1111111"
const UserEmail = "arnaud@artisons.me"
const AdminEmail = "hello@artisons.me"
const OtpNotMatching = "otp@artisons.me"
const AdminBlockedEmail = "blocked@artisons.me"
const IDDoesNotExist = 9999999
const EmailDoesNotExist = "idontexist@artisons.me"
const DoesNotExist = "idontexist"
const SeoKey = "test"
const SeoTitle = "The social networks are evil."
const SeoURL = "/test.html"
const SeoDescription = "Buh the social networks."
const FilterColorKey = "colors"
const FilterColorLabel = "colors"
const FilterSizeKey = "sizes"
const FilterToDeleteKey = "delete"
const Tag = "mens"
const Branch1 = "clothes"
const Branch2 = "shoes"
const Tag2 = "womens"
const TagToDelete = "children"
const UA = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
const Delivery = "colissimo"
const OrderID = "ORD1"
const OrderID2 = "ORD2"
const OrderFirstName = "Arnaud"
const OrderTotal = 105.5

var ArticleSlug string = stringutil.Slugify(db.Escape(ArticleTitle))

func AddToCart(ctx context.Context, cid, pid string, qty string) {
	db.Redis.HSet(ctx, "cart:"+cid, pid, qty).Result()
	db.Redis.Expire(ctx, "cart:"+cid, conf.CartDuration)
}

func Quantity(ctx context.Context, cid, pid string) string {
	qty, _ := db.Redis.HGet(ctx, "cart:"+cid, pid).Result()

	return qty
}

func TTL(ctx context.Context, key string) time.Duration {
	ttl, _ := db.Redis.TTL(ctx, key).Result()
	return ttl
}

func init() {
	ctx := Context()

	lines := db.ParseData(ctx, conf.WorkingSpace+"web/redis/unit.redis")
	pipe := db.Redis.Pipeline()

	for _, line := range lines {
		pipe.Do(ctx, line...)
		// log.Println(line)
		// db.Redis.Do(ctx, line).Result()
	}

	cmds, err := pipe.Exec(ctx)

	for _, cmd := range cmds {
		if cmd.Err() != nil {
			log.Println(cmd.String(), cmd.Err())
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	// db.Redis.HSet(ctx, "filter:"+FilterToDeleteKey,
	// 	"key", FilterToDeleteKey,
	// 	"editable", "1",
	// 	"label", "Sizes",
	// 	"values", "S;M;L",
	// 	"updated_at", time.Now().Unix(),
	// )

	// db.Redis.ZAdd(ctx, "filters", redis.Z{
	// 	Score:  float64(now.Unix()),
	// 	Member: FilterColorKey,
	// }, redis.Z{
	// 	Score:  float64(now.Unix()),
	// 	Member: FilterSizeKey,
	// }, redis.Z{
	// 	Score:  float64(now.Unix()),
	// 	Member: FilterToDeleteKey,
	// })

	// db.Redis.ZAdd(ctx, "filters:active", redis.Z{
	// 	Score:  float64(1),
	// 	Member: FilterColorKey,
	// })

	// db.Redis.ZAdd(ctx, fmt.Sprintf("wish:%d", UserID1),
	// 	redis.Z{
	// 		Member: "PDT1",
	// 		Score:  float64(now.Unix()),
	// 	},
	// 	redis.Z{
	// 		Member: "PDT2",
	// 		Score:  float64(now.Unix()),
	// 	},
	// )

	// db.Redis.Del(ctx, "tags")
	// db.Redis.Del(ctx, "tags:root")
	// db.Redis.HSet(ctx, "tag:"+Tag, "key", Tag, "image", "tags/1.jpeg", "label", "Mens", "order", "1", "children", "clothes;shoes")
	// db.Redis.HSet(ctx, "tag:"+Tag2, "key", Tag2, "image", "tags/2.jpeg", "label", "Womens", "order", "2", "children", "clothes;shoes")
	// db.Redis.HSet(ctx, "tag:"+TagToDelete, "key", TagToDelete, "image", "tags/3.jpeg", "label", "Children", "order", "2")
	// db.Redis.HSet(ctx, "tag:"+Branch1, "key", Branch1, "image", "tags/3.jpeg", "label", "Clothes", "order", "2")
	// db.Redis.HSet(ctx, "tag:"+Branch2, "key", Branch2, "image", "tags/3.jpeg", "label", "Shoes", "order", "2")

	// db.Redis.ZAdd(ctx, "tags", redis.Z{
	// 	Score:  float64(1),
	// 	Member: Tag,
	// }, redis.Z{
	// 	Score:  float64(2),
	// 	Member: Tag2,
	// }, redis.Z{
	// 	Score:  float64(2),
	// 	Member: Branch1,
	// }, redis.Z{
	// 	Score:  float64(2),
	// 	Member: Branch2,
	// })

	// db.Redis.ZAdd(ctx, "tags:root", redis.Z{
	// 	Score:  float64(1),
	// 	Member: Tag,
	// }, redis.Z{
	// 	Score:  float64(2),
	// 	Member: Tag2,
	// })

	// db.Redis.ZAdd(ctx, "deliveries", redis.Z{
	// 	Score:  1,
	// 	Member: "colissimo",
	// })

	// db.Redis.HSet(ctx, fmt.Sprintf("user:%d", UserID1),
	// 	"id", UserID1,
	// 	"email", UserEmail,
	// 	"created_at", now.Unix(),
	// 	"updated_at", now.Unix(),
	// 	"type", "user",
	// 	"role", "user",
	// )

	// db.Redis.HSet(ctx, "user:2",
	// 	"id", "2",
	// 	"email", AdminEmail,
	// 	"created_at", now.Unix(),
	// 	"updated_at", now.Unix(),
	// 	"type", "user",
	// 	"role", "admin",
	// )

	// db.Redis.HSet(ctx, "session:"+UserSID, "id", UserSID, "uid", UserID1, "device", UA, "type", "session")
	// db.Redis.Expire(ctx, "session:"+UserSID, conf.SessionDuration)
	// db.Redis.HSet(ctx, "session:"+UserSIDSignedIn, "id", UserSIDSignedIn, "uid", "3", "type", "session")
	// db.Redis.Expire(ctx, "session:"+UserSIDSignedIn, conf.SessionDuration)

	// referers := []string{"Google", "Unknown", "Yandex", "DuckDuckGo"}
	// products := []string{"PDT1", "PDT2", "PDT3", "PDT4", "PDT5"}
	// browsers := []string{"Chrome", "Safari", "Firefox", "Edge"}
	// systems := []string{"Windows", "Android", "iOS", "Linux"}
	// urls := []string{"index.html", "/super-article-du-blog.html", "/PDT2-sweat-a-capuche-uniforme.html", "/cgv.html", "/panier.html", "/coucou.html"}
	// pipe := db.Redis.Pipeline()

	// now := time.Now()
	// for i := 0; i < 50; i++ {
	// 	i := rand.Intn(30)
	// 	score := now.AddDate(0, 0, -i)

	// 		log.SetFlags(0)
	// log.Printf("ZINCRBY demo:stats:pageviews:20060102 1 %s", slug)
	// log.Printf("ZINCRBY demo:stats:products:most:20060102 1 %s", products[prand])
	// log.Printf("ZINCRBY demo:stats:products:shared:20060102 1 %s", products[prand])
	// log.Printf("ZINCRBY demo:stats:browsers:20060102 1 %s", browsers[brand])
	// log.Printf("ZINCRBY demo:stats:referers:20060102 1 %s", referers[rrand])
	// log.Printf("ZINCRBY demo:stats:referers:20060102 1 %s", referers[rrand])
	// log.Printf("ZINCRBY demo:stats:systems:20060102 1 %s", systems[srand])
	// log.Printf("SET demo:stats:visits:20060102 %d 0", visits)
	// log.Printf("SET demo:stats:orders:revenues:20060102 1 %d 0", amount)
	// log.Printf("SET demo:stats:orders:count:20060102 1 %d 0", count)
	// log.Printf("SET demo:stats:visits:unique:20060102 1 %d 0", uniques)
	// log.Printf("SET demo:stats:pageviews:all:20060102 1 %d 0", visits*2)
	// }
	// 	urli := rand.Intn(len(urls))
	// 	slug := urls[urli]

	// 	visits := rand.Intn(1000) + 100
	// 	amount := rand.Intn(5000) + 100
	// 	count := rand.Intn(10) + 1
	// 	uniques := rand.Intn(visits)
	// 	rrand := rand.Intn(len(referers))
	// 	brand := rand.Intn(len(browsers))
	// 	srand := rand.Intn(len(systems))
	// 	prand := rand.Intn(len(products))

	// 	pipe.ZIncrBy(ctx, "demo:stats:pageviews:"+score.Format("20060102"), 1, slug)
	// 	pipe.ZIncrBy(ctx, "demo:stats:products:most:"+score.Format("20060102"), 1, products[prand])
	// 	pipe.ZIncrBy(ctx, "demo:stats:products:shared:"+score.Format("20060102"), 1, products[prand])
	// 	pipe.ZIncrBy(ctx, "demo:stats:browsers:"+score.Format("20060102"), 1, browsers[brand])
	// 	pipe.ZIncrBy(ctx, "demo:stats:referers:"+score.Format("20060102"), 1, referers[rrand])
	// 	pipe.ZIncrBy(ctx, "demo:stats:systems:"+score.Format("20060102"), 1, systems[srand])
	// 	pipe.Set(ctx, "demo:stats:visits:"+score.Format("20060102"), visits, 0)
	// 	pipe.Set(ctx, "demo:stats:orders:revenues:"+score.Format("20060102"), amount, 0)
	// 	pipe.Set(ctx, "demo:stats:orders:count:"+score.Format("20060102"), count, 0)
	// 	pipe.Set(ctx, "demo:stats:visits:unique:"+score.Format("20060102"), uniques, 0)
	// 	pipe.Set(ctx, "demo:stats:pageviews:all:"+score.Format("20060102"), visits*2, 0)
	// }

	// log.SetFlags(0)
	// log.Printf("ZINCRBY demo:stats:pageviews:20060102 1 %s", slug)
	// log.Printf("ZINCRBY demo:stats:products:most:20060102 1 %s", products[prand])
	// log.Printf("ZINCRBY demo:stats:products:shared:20060102 1 %s", products[prand])
	// log.Printf("ZINCRBY demo:stats:browsers:20060102 1 %s", browsers[brand])
	// log.Printf("ZINCRBY demo:stats:referers:20060102 1 %s", referers[rrand])
	// log.Printf("ZINCRBY demo:stats:referers:20060102 1 %s", referers[rrand])
	// log.Printf("ZINCRBY demo:stats:systems:20060102 1 %s", systems[srand])
	// log.Printf("SET demo:stats:visits:20060102 %d 0", visits)
	// log.Printf("SET demo:stats:orders:revenues:20060102 1 %d 0", amount)
	// log.Printf("SET demo:stats:orders:count:20060102 1 %d 0", count)
	// log.Printf("SET demo:stats:visits:unique:20060102 1 %d 0", uniques)
	// log.Printf("SET demo:stats:pageviews:all:20060102 1 %d 0", visits*2)

	// if _, err := pipe.Exec(ctx); err != nil {
	// 	log.Fatalln(err)
	// }
}

func Context() context.Context {
	var ctx context.Context = context.WithValue(context.Background(), middleware.RequestIDKey, fmt.Sprintf("%d", time.Now().UnixMilli()))
	ctx = context.WithValue(ctx, contexts.Locale, language.English)

	rid, _ := stringutil.Random()
	ctx = context.WithValue(ctx, middleware.RequestIDKey, rid)
	ctx = context.WithValue(ctx, contexts.Demo, true)

	return context.WithValue(ctx, contexts.Device, fmt.Sprintf("%d", time.Now().UnixMilli()))
}
