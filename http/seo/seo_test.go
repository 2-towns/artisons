package seo

import (
	"artisons/conf"
	"artisons/db"
	"artisons/tests"
	"log"
	"testing"
	"time"
)

var content = Content{
	Key:         "test",
	URL:         "/test.html",
	Title:       "The social networks are evil.",
	Description: "Buh the social networks.",
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
