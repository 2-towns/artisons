package seo

import (
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/router"
	"gifthub/tests"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var content = Content{
	Key:         "test",
	URL:         "/test.html",
	Title:       "The social networks are evil.",
	Description: "Buh the social networks.",
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("coucou"))
}

func init() {
	ctx := tests.Context()
	pipe := db.Redis.Pipeline()
	now := time.Now().Unix()

	pipe.HSet(ctx, "seo:"+content.Key,
		"title", db.Escape(content.Title),
		"description", db.Escape(content.Description),
		"url", db.Escape(content.URL),
		"key", content.Key,
		"lang", conf.DefaultLocale.String(),
		"updated_at", now,
	)

	if _, err := pipe.Exec(ctx); err != nil {
		log.Panicln(err)
	}

	URLs[content.Key] = content

	router.R.Get(content.URL, healthCheckHandler)
}

func TestValidateReturnsErrorWhenKeyIsEmpty(t *testing.T) {
	ctx := tests.Context()

	c := content
	c.Key = ""

	if err := c.Validate(ctx); err == nil || err.Error() != "input:key" {
		t.Fatalf(`c.Validate(ctx) = %v, want not "input:key"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenTitleIsEmpty(t *testing.T) {
	ctx := tests.Context()

	c := content
	c.Title = ""

	if err := c.Validate(ctx); err == nil || err.Error() != "input:title" {
		t.Fatalf(`c.Validate(ctx) = %v, want not "input:title"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenURLIsEmpty(t *testing.T) {
	ctx := tests.Context()

	c := content
	c.URL = ""

	if err := c.Validate(ctx); err == nil || err.Error() != "input:url" {
		t.Fatalf(`c.Validate(ctx) = %v, want not "input:url"`, err.Error())
	}
}

func TestGetReturnsContentWhenSeoExists(t *testing.T) {
	ctx := tests.Context()

	c, err := Find(ctx, "test")

	if err != nil {
		t.Fatalf(`Find(ctx,"test",conf.DefaultLocale.String()) = _, %s, want _, nil`, err.Error())
	}

	if c.Key == "" {
		t.Fatalf(`c.Key = %s, want not empty`, c.Key)
	}

	if c.Title == "" {
		t.Fatalf(`c.Title = %s, want not empty`, c.Title)
	}

	if c.Description == "" {
		t.Fatalf(`c.Description = %s, want not empty`, c.Description)
	}

	if c.URL == "" {
		t.Fatalf(`c.URL = %s, want not empty`, c.URL)
	}
}

func TestGetReturnsErrorContentWhenSeoDoesNotExist(t *testing.T) {
	ctx := tests.Context()

	if _, err := Find(ctx, "jenexistepas"); err == nil || err.Error() != "the data is not found" {
		t.Fatalf(`Find(ctx,"jenexistepas",conf.DefaultLocale.String()) = _, %s, want _, "the data is not found"`, err)
	}
}

func TestSaveReturnsNoErrorWhenDataAreOk(t *testing.T) {
	svr := httptest.NewServer(router.R)
	defer svr.Close()

	res, err := http.Get(svr.URL + "/test.html")

	if err != nil {
		t.Fatalf(`http.Get(svr.URL + "/test.html") = _, %v, want _, nil`, err.Error())
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf(`res.StatusCode = %d, want %d`, res.StatusCode, http.StatusOK)
	}

	res, err = http.Get(svr.URL + "/test-modified.html")

	if err != nil {
		t.Fatalf(`http.Get(svr.URL + "/test-modified.html") = _, %v, want _, nil`, err.Error())
	}

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf(`res.StatusCode = %d, want %d`, res.StatusCode, http.StatusNotFound)
	}

	ctx := tests.Context()
	c := Content{
		Key:         "test",
		URL:         "/test-modified.html",
		Title:       "The social networks are evil.",
		Description: "Buh the social networks.",
	}

	if _, err := c.Save(ctx); err != nil {
		t.Fatalf(`c.Save(ctx) = %s, want nil`, err.Error())
	}

	res, err = http.Get(svr.URL + "/test.html")

	if err != nil {
		t.Fatalf(`http.Get(svr.URL + "/testr.html") = _, %v, want _, nil`, err.Error())
	}

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf(`res.StatusCode = %d, want %d`, res.StatusCode, http.StatusNotFound)
	}

	res, err = http.Get(svr.URL + "/test-modified.html")

	if err != nil {
		t.Fatalf(`http.Get(svr.URL + "/test-modified.html") = _, %v, want _, nil`, err.Error())
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf(`res.StatusCode = %d, want %d`, res.StatusCode, http.StatusOK)
	}

}
