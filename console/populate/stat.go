package populate

import (
	"context"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

func stats(ctx context.Context, pipe redis.Pipeliner) {
	now := time.Now()

	referers := []string{"Google", "Unknown", "Yandex", "DuckDuckGo"}
	products := []string{"PDT1", "PDT2", "PDT3", "PDT4", "PDT5"}
	browsers := []string{"Chrome", "Safari", "Firefox", "Edge"}
	systems := []string{"Windows", "Android", "iOS", "Linux"}
	urls := []string{"index.html", "/super-article-du-blog.html", "/PDT2-sweat-a-capuche-uniforme.html", "/cgv.html", "/panier.html", "/coucou.html"}

	for i := 0; i < 5000; i++ {
		i := rand.Intn(30)
		score := now.AddDate(0, 0, -i)

		urli := rand.Intn(len(urls))
		slug := urls[urli]

		visits := rand.Intn(1000) + 100
		amount := rand.Intn(5000) + 100
		count := rand.Intn(10) + 1
		uniques := rand.Intn(visits)
		rrand := rand.Intn(len(referers))
		brand := rand.Intn(len(browsers))
		srand := rand.Intn(len(systems))
		prand := rand.Intn(len(products))

		pipe.ZIncrBy(ctx, "demo:stats:pageviews:"+score.Format("20060102"), 1, slug)
		pipe.ZIncrBy(ctx, "demo:stats:products:most:"+score.Format("20060102"), 1, products[prand])
		pipe.ZIncrBy(ctx, "demo:stats:products:shared:"+score.Format("20060102"), 1, products[prand])
		pipe.ZIncrBy(ctx, "demo:stats:browsers:"+score.Format("20060102"), 1, browsers[brand])
		pipe.ZIncrBy(ctx, "demo:stats:referers:"+score.Format("20060102"), 1, referers[rrand])
		pipe.ZIncrBy(ctx, "demo:stats:systems:"+score.Format("20060102"), 1, systems[srand])
		pipe.Set(ctx, "demo:stats:visits:"+score.Format("20060102"), visits, 0)
		pipe.Set(ctx, "demo:stats:orders:revenues:"+score.Format("20060102"), amount, 0)
		pipe.Set(ctx, "demo:stats:orders:count:"+score.Format("20060102"), count, 0)
		pipe.Set(ctx, "demo:stats:visits:unique:"+score.Format("20060102"), uniques, 0)
		pipe.Set(ctx, "demo:stats:pageviews:all:"+score.Format("20060102"), visits*2, 0)
	}
}
