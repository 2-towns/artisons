package seo

import (
	"artisons/tests"
	"testing"
)

var content = Content{
	Key:         tests.SeoKey,
	URL:         tests.SeoURL,
	Title:       tests.SeoTitle,
	Description: tests.SeoDescription,
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

	c, err := Find(ctx, content.Key)

	if err != nil {
		t.Fatalf(`Find(ctx, content.Key, conf.DefaultLocale.String()) = _, %s, want _, nil`, err.Error())
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

	if _, err := Find(ctx, tests.DoesNotExist); err == nil || err.Error() != "the data is not found" {
		t.Fatalf(`Find(ctx, tests.DoesNotExist, conf.DefaultLocale.String()) = _, %s, want _, "the data is not found"`, err)
	}
}
