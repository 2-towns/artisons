package stats

import (
	"artisons/http/contexts"
	"artisons/tests"
	"artisons/users"
	"context"
	"testing"
)

func TestGetAll(t *testing.T) {
	c := tests.Context()
	c = context.WithValue(c, contexts.User, users.User{Demo: true})

	data, err := GetAll(c, 14)

	if err != nil {
		t.Fatalf(`err = %s, want nil`, err.Error())
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

// func TestMostValuesReturnsDataWhenSuccess(t *testing.T) {
// 	c := tests.Context()

// 	data, err := MostValues(c, 14)

// 	if err != nil {
// 		t.Fatalf(`MostValues(c, 14) = %s, want nil`, err.Error())
// 	}

// 	if len(data) != 6 {
// 		t.Fatalf(`len(data) = %v, want 6`, len(data))
// 	}

// 	pageviews := data[0]

// 	if len(pageviews) == 0 {
// 		t.Fatalf(`len(pageviews) = %v, want > 0`, len(pageviews))
// 	}

// 	item := pageviews[0]

// 	if item.Key == "" {
// 		t.Fatalf(`item.Key = %v, want not empty`, item.Key)
// 	}

// 	if item.URL == "" {
// 		t.Fatalf(`item.URL = %v, want not empty`, item.URL)
// 	}

// 	if item.Percent == 0 {
// 		t.Fatalf(`item.Percent = %v, want > 0`, item.Percent)
// 	}

// 	if item.Value == 0 {
// 		t.Fatalf(`item.Value = %v, want > 0`, item.Value)
// 	}
// }

// func TestVisitIncrementsDataWhenSuccess(t *testing.T) {
// 	ctx := tests.Context()
// 	ua := useragent.Parse(tests.UA)

// 	bef, err := GetAll(ctx, 1)
// 	if err != nil {
// 		t.Fatalf(`Visit(c, ua, VisitData{URL: "/index.html", Referer: ""}) = %s, want nil`, err.Error())
// 	}

// 	mbeft, err := MostValues(ctx, 1)

// 	err = Visit(ctx, ua, VisitData{URL: "/index.html", Referer: ""})
// 	if err != nil {
// 		t.Fatalf(`Visit(c, ua, VisitData{URL: "/index.html", Referer: ""}) = %s, want nil`, err.Error())
// 	}

// 	aft, err := GetAll(ctx, 1)
// 	if err != nil {
// 		t.Fatalf(`Visit(c, ua, VisitData{URL: "/index.html", Referer: ""}) = %s, want nil`, err.Error())
// 	}

// 	maft, err := MostValues(ctx, 1)
// 	if err != nil {
// 		t.Fatalf(`Visit(c, ua, VisitData{URL: "/index.html", Referer: ""}) = %s, want nil`, err.Error())
// 	}

// 	log.Println(mbeft[0])
// 	log.Println(maft[0])

// 	// Visits
// 	if bef[1].Sum < aft[1].Sum {
// 		t.Fatalf(`Visit(c, ua, VisitData{URL: "/index.html", Referer: ""}) = %s, want nil`, err.Error())
// 	}

// 	// Pages views
// 	if bef[2].Sum >= aft[2].Sum {
// 		t.Fatalf(`Visit(c, ua, VisitData{URL: "/index.html", Referer: ""}) = %s, want nil`, err.Error())
// 	}
// }

// func TestVisitIncrementsRefererWhenRefererIsSet(t *testing.T) {
// 	ctx := tests.Context()
// 	ua := useragent.Parse(tests.UA)
// 	now := time.Now().Format("20060102")
// 	visits, _ := db.Redis.Get(ctx, "stats:visits:"+now).Result()
// 	ivisits, _ := strconv.ParseInt(visits, 10, 64)
// 	referers, _ := db.Redis.ZScore(ctx, "stats:referers:"+now, "Google").Result()

// 	err := Visit(ctx, ua, VisitData{URL: "/index.html", Referer: "http://www.google.com/search"})
// 	if err != nil {
// 		t.Fatalf(`Visit(c, ua, VisitData{URL: "/index.html", Referer: "http://www.google.com/search"}) = %s, want nil`, err.Error())
// 	}

// 	v, _ := db.Redis.Get(ctx, "stats:visits:"+now).Result()
// 	if fmt.Sprintf("%d", ivisits+1) != v {
// 		t.Fatalf(`visits = %s, want %d`, v, ivisits+1)
// 	}

// 	r, _ := db.Redis.ZScore(ctx, "stats:referers:"+now, "Google").Result()
// 	if referers+1 != r {
// 		t.Fatalf(`referers = %f, want %f`, r, referers+1)
// 	}
// }

// func TestVisitDoesNotIncrementUniqueWhenAlreadyVisited(t *testing.T) {
// 	ctx := tests.Context()
// 	ua := useragent.Parse(tests.UA)
// 	now := time.Now().Format("20060102")
// 	db.Redis.SAdd(ctx, "stats:visits:members:"+now, ctx.Value(contexts.Device))
// 	uniques, _ := db.Redis.Get(ctx, "stats:visits:unique:"+now).Result()

// 	err := Visit(ctx, ua, VisitData{URL: "/index.html", Referer: ""})
// 	if err != nil {
// 		t.Fatalf(`Visit(c, ua, VisitData{URL: "/index.html", Referer: ""}) = %s, want nil`, err.Error())
// 	}

// 	u, _ := db.Redis.Get(ctx, "stats:visits:unique:"+now).Result()
// 	if uniques != u {
// 		t.Fatalf(`uniques = %s, want %s`, u, uniques)
// 	}
// }

// func TestOrderIncrementsDataWhenValid(t *testing.T) {
// 	ctx := tests.Context()
// 	oid, _ := stringutil.Random()
// 	pid, _ := stringutil.Random()
// 	q := map[string]int{
// 		pid: 2,
// 	}
// 	total := 102.4
// 	now := time.Now().Format("20060102")

// 	count, _ := db.Redis.Get(ctx, "stats:orders:count:"+now).Result()
// 	icount, _ := strconv.ParseInt(count, 10, 64)

// 	revenues, _ := db.Redis.Get(ctx, "stats:orders:revenues:"+now).Result()
// 	frevenues, _ := strconv.ParseFloat(revenues, 64)

// 	err := Order(ctx, oid, q, total)
// 	if err != nil {
// 		t.Fatalf(`Order(ctx, oid, q, total) = %s, want nil`, err.Error())
// 	}

// 	p, _ := db.Redis.ZScore(ctx, "stats:products:most:"+now, pid).Result()
// 	if p != 2 {
// 		t.Fatalf(`%s = %f, want 2`, pid, p)
// 	}

// 	c, _ := db.Redis.Get(ctx, "stats:orders:count:"+now).Result()
// 	if fmt.Sprintf("%d", icount+1) != c {
// 		t.Fatalf(`count = %s, want %d`, c, icount+1)
// 	}

// 	r, _ := db.Redis.Get(ctx, "stats:orders:revenues:"+now).Result()
// 	fr, _ := strconv.ParseFloat(r, 64)
// 	if fmt.Sprintf("%f", fr) != fmt.Sprintf("%f", frevenues+total) {
// 		t.Fatalf(`revenues = %f, want %f`, fr, frevenues+total)
// 	}
// }
