package stats

import (
	"fmt"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/string/stringutil"
	"artisons/tests"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/mileusna/useragent"
)

func init() {
	ctx := tests.Context()
	referers := []string{"Google", "Unknown", "Yandex", "DuckDuckGo"}
	products := []string{"PDT1", "PDT2", "PDT3", "PDT4", "PDT5"}
	browsers := []string{"Chrome", "Safari", "Firefox", "Edge"}
	systems := []string{"Windows", "Android", "iOS", "Linux"}
	urls := []string{"index.html", "/super-article-du-blog.html", "/PDT2-sweat-a-capuche-uniforme.html", "/cgv.html", "/panier.html", "/coucou.html"}
	now := time.Now()
	pipe := db.Redis.Pipeline()

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

	if _, err := pipe.Exec(ctx); err != nil {
		log.Fatalln(err)
	}

}

func TestGellAllReturnsDataWhenSuccess(t *testing.T) {
	c := tests.Context()

	data, err := GetAll(c, 14)

	if err != nil {
		t.Fatalf(`GetAll(c, 14) = %s, want nil`, err.Error())
	}

	if len(data) != 6 {
		t.Fatalf(`len(data) = %v, want 6`, len(data))
	}

	v := data[0]
	if len(v.Value) != 14 {
		t.Fatalf(`len(data) = %v, want 14`, len(v.Value))
	}

	if v.Sum == 0 {
		t.Fatalf(`len(data) = %d, want > 0`, v.Sum)
	}
}

func TestMostValuesReturnsDataWhenSuccess(t *testing.T) {
	c := tests.Context()

	data, err := MostValues(c, 14)

	if err != nil {
		t.Fatalf(`MostValues(c, 14) = %s, want nil`, err.Error())
	}

	if len(data) != 6 {
		t.Fatalf(`len(data) = %v, want 6`, len(data))
	}

	pageviews := data[0]

	if len(pageviews) == 0 {
		t.Fatalf(`len(pageviews) = %v, want > 0`, len(pageviews))
	}

	item := pageviews[0]

	if item.Key == "" {
		t.Fatalf(`item.Key = %v, want not empty`, item.Key)
	}

	if item.URL == "" {
		t.Fatalf(`item.URL = %v, want not empty`, item.URL)
	}

	if item.Percent == 0 {
		t.Fatalf(`item.Percent = %v, want > 0`, item.Percent)
	}

	if item.Value == 0 {
		t.Fatalf(`item.Value = %v, want > 0`, item.Value)
	}
}

func TestVisitIncrementsDataWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	ua := useragent.Parse("Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	now := time.Now().Format("20060102")
	visits, _ := db.Redis.Get(ctx, "stats:visits:"+now).Result()
	browsers, _ := db.Redis.ZScore(ctx, "stats:browsers:"+now, ua.Name).Result()
	pageviews, _ := db.Redis.ZScore(ctx, "stats:pageviews:"+now, "/index.html").Result()
	allpageviews, _ := db.Redis.Get(ctx, "stats:pageviews:all:"+now).Result()
	iallpageviews, _ := strconv.ParseInt(allpageviews, 10, 64)
	unique, _ := db.Redis.Get(ctx, "stats:visits:unique:"+now).Result()
	iunique, _ := strconv.ParseInt(unique, 10, 64)

	err := Visit(ctx, ua, VisitData{URL: "/index.html", Referer: ""})
	if err != nil {
		t.Fatalf(`Visit(c, ua, VisitData{URL: "/index.html", Referer: ""}) = %s, want nil`, err.Error())
	}

	v, _ := db.Redis.Get(ctx, "stats:visits:"+now).Result()
	if visits != v {
		t.Fatalf(`visits = %s, want %s`, v, visits)
	}

	b, _ := db.Redis.ZScore(ctx, "stats:browsers:"+now, ua.Name).Result()
	if browsers+1 != b {
		t.Fatalf(`browsers = %f, want %f`, b, browsers+1)
	}

	p, _ := db.Redis.ZScore(ctx, "stats:pageviews:"+now, "/index.html").Result()
	if pageviews+1 != p {
		t.Fatalf(`pageviews = %f, want %f`, p, pageviews+1)
	}

	a, _ := db.Redis.Get(ctx, "stats:pageviews:all:"+now).Result()
	if fmt.Sprintf("%d", iallpageviews+1) != a {
		t.Fatalf(`pageviews:all = %s, want %d`, a, iallpageviews+1)
	}

	u, _ := db.Redis.Get(ctx, "stats:visits:unique:"+now).Result()
	if fmt.Sprintf("%d", iunique+1) != u {
		t.Fatalf(`pageviews:all = %s, want %d`, u, iunique+1)
	}

	exists, _ := db.Redis.SIsMember(ctx, "stats:visits:members:"+now, ctx.Value(contexts.Cart)).Result()
	if !exists {
		t.Fatalf(`exists = %v, want true`, exists)
	}

	ttl, _ := db.Redis.TTL(ctx, "stats:visits:members:"+now).Result()
	if ttl != time.Hour*24 {
		t.Fatalf(`ttl = %d, want %d`, ttl, time.Hour*24)
	}
}

func TestVisitIncrementsRefererWhenRefererIsSet(t *testing.T) {
	ctx := tests.Context()
	ua := useragent.Parse("Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	now := time.Now().Format("20060102")
	visits, _ := db.Redis.Get(ctx, "stats:visits:"+now).Result()
	ivisits, _ := strconv.ParseInt(visits, 10, 64)
	referers, _ := db.Redis.ZScore(ctx, "stats:referers:"+now, "Google").Result()

	err := Visit(ctx, ua, VisitData{URL: "/index.html", Referer: "http://www.google.com/search"})
	if err != nil {
		t.Fatalf(`Visit(c, ua, VisitData{URL: "/index.html", Referer: "http://www.google.com/search"}) = %s, want nil`, err.Error())
	}

	v, _ := db.Redis.Get(ctx, "stats:visits:"+now).Result()
	if fmt.Sprintf("%d", ivisits+1) != v {
		t.Fatalf(`visits = %s, want %d`, v, ivisits+1)
	}

	r, _ := db.Redis.ZScore(ctx, "stats:referers:"+now, "Google").Result()
	if referers+1 != r {
		t.Fatalf(`referers = %f, want %f`, r, referers+1)
	}
}

func TestVisitDoesNotIncrementUniqueWhenAlreadyVisited(t *testing.T) {
	ctx := tests.Context()
	ua := useragent.Parse("Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	now := time.Now().Format("20060102")
	db.Redis.SAdd(ctx, "stats:visits:members:"+now, ctx.Value(contexts.Cart))
	uniques, _ := db.Redis.Get(ctx, "stats:visits:unique:"+now).Result()

	err := Visit(ctx, ua, VisitData{URL: "/index.html", Referer: ""})
	if err != nil {
		t.Fatalf(`Visit(c, ua, VisitData{URL: "/index.html", Referer: ""}) = %s, want nil`, err.Error())
	}

	u, _ := db.Redis.Get(ctx, "stats:visits:unique:"+now).Result()
	if uniques != u {
		t.Fatalf(`uniques = %s, want %s`, u, uniques)
	}
}

func TestOrderIncrementsDataWhenValid(t *testing.T) {
	ctx := tests.Context()
	oid, _ := stringutil.Random()
	pid, _ := stringutil.Random()
	q := map[string]int{
		pid: 2,
	}
	total := 102.4
	now := time.Now().Format("20060102")

	count, _ := db.Redis.Get(ctx, "stats:orders:count:"+now).Result()
	icount, _ := strconv.ParseInt(count, 10, 64)

	revenues, _ := db.Redis.Get(ctx, "stats:orders:revenues:"+now).Result()
	frevenues, _ := strconv.ParseFloat(revenues, 64)

	err := Order(ctx, oid, q, total)
	if err != nil {
		t.Fatalf(`Order(ctx, oid, q, total) = %s, want nil`, err.Error())
	}

	p, _ := db.Redis.ZScore(ctx, "stats:products:most:"+now, pid).Result()
	if p != 2 {
		t.Fatalf(`%s = %f, want 2`, pid, p)
	}

	c, _ := db.Redis.Get(ctx, "stats:orders:count:"+now).Result()
	if fmt.Sprintf("%d", icount+1) != c {
		t.Fatalf(`count = %s, want %d`, c, icount+1)
	}

	r, _ := db.Redis.Get(ctx, "stats:orders:revenues:"+now).Result()
	fr, _ := strconv.ParseFloat(r, 64)
	if fmt.Sprintf("%f", fr) != fmt.Sprintf("%f", frevenues+total) {
		t.Fatalf(`revenues = %f, want %f`, fr, frevenues+total)
	}
}

func TestMiddlewareUpdatesStatWhenNotAdmin(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/index.html", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := tests.Context()
	now := time.Now().Format("20060102")
	pageviews, _ := db.Redis.ZScore(ctx, "stats:pageviews:"+now, "/index.html").Result()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	handler := Middleware(testHandler)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf(`status = %d, want %d`, status, http.StatusOK)
	}

	time.Sleep(time.Millisecond * 10)

	p, _ := db.Redis.ZScore(ctx, "stats:pageviews:"+now, "/index.html").Result()
	if pageviews+1 != p {
		t.Fatalf(`pageviews = %f, want %f`, p, pageviews+1)
	}
}
