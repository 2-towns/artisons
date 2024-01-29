package blog

import (
	"artisons/conf"
	"artisons/db"
	"artisons/string/stringutil"
	"artisons/tests"
	"os"
	"path"
	"testing"
	"time"
)

var article Article = Article{
	Title:       "Mangez de l'ail !",
	Description: "C'est un antiseptique.",
	Slug:        db.Escape("Mangez de l'ail !"),
	Image:       path.Join(conf.WorkingSpace, "web", "tmp", "hello"),
	Status:      "online",
}

func init() {
	ctx := tests.Context()
	now := time.Now()

	db.Redis.Del(ctx, "blog")
	db.Redis.HSet(ctx, "blog:99",
		"id", "99",
		"title", article.Title,
		"status", article.Status,
		"slug", stringutil.Slugify(db.Escape(article.Title)),
		"lang", conf.DefaultLocale.String(),
		"description", article.Description,
		"image", path.Join(conf.WorkingSpace, "web", "images", "blog", "99.jpeg"),
		"updated_at", now.Unix(),
		"type", "blog",
	)

	db.Redis.HSet(ctx, "bids", db.Escape(article.Title), "99")

	db.Redis.HSet(ctx, "blog:98",
		"id", "98",
		"title", article.Title,
		"status", article.Status,
		"lang", conf.DefaultLocale.String(),
		"slug", stringutil.Slugify(article.Title+"-bis"),
		"description", article.Description,
		"image", path.Join(conf.WorkingSpace, "web", "images", "blog", "99.jpeg"),
		"type", "blog",
		"updated_at", now.Unix(),
	)

	db.Redis.HSet(ctx, "blog:97",
		"id", "97",
		"title", "Terms of Sales",
		"status", "online",
		"lang", conf.DefaultLocale.String(),
		"description", "Hello !",
		"type", "cms",
		"updated_at", now.Unix(),
	)
}

func TestValidateReturnsErrorWhenTitleIsEmpty(t *testing.T) {
	c := tests.Context()

	a := article
	a.Title = ""

	if err := a.Validate(c); err == nil || err.Error() != "input:title" {
		t.Fatalf(`a.Validate(c) = %v, want "input:title"`, err.Error())
	}
}

func TestValidReturnsErrorWhenDescriptionIsEmpty(t *testing.T) {
	c := tests.Context()

	a := article
	a.Description = ""

	if err := a.Validate(c); err == nil || err.Error() != "input:description" {
		t.Fatalf(`a.Validate(c) = %v, want "input:description"`, err.Error())
	}
}

func TestValidReturnsErrorWhenSlugIsEmpty(t *testing.T) {
	c := tests.Context()

	a := article
	a.Slug = ""

	if err := a.Validate(c); err == nil || err.Error() != "input:slug" {
		t.Fatalf(`a.Validate(c) = %v, want "input:slug"`, err.Error())
	}
}

func TestValidateReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	os.Create(path.Join(conf.WorkingSpace, "web", "tmp", "hello"))

	if err := article.Validate(c); err != nil {
		t.Fatalf(`a.Validate(c) = %v, want nil`, err.Error())
	}
}

func TestSearchReturnsArticlesWhenTitleIsFound(t *testing.T) {
	c := tests.Context()
	a, err := Search(c, Query{Keywords: "Mangez de l'ail"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: "Mangez de l'ail"}) = %v, want nil`, err.Error())
	}

	if a.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, a.Total)
	}

	if len(a.Articles) == 0 {
		t.Fatalf(`len(p.Articles) = %d, want > 0`, len(a.Articles))
	}

	if a.Articles[0].ID == 99 {
		t.Fatalf(`p[0].ID = %d, want = 99`, a.Articles[0].ID)
	}
}

func TestSearchReturnsArticlesWhenAtLeastOneMatching(t *testing.T) {
	c := tests.Context()
	a, err := Search(c, Query{Keywords: "hello ail"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: ""hello ail"}) = %v, want nil`, err.Error())
	}

	if a.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, a.Total)
	}

	if len(a.Articles) == 0 {
		t.Fatalf(`len(p.Articles) = %d, want > 0`, len(a.Articles))
	}

	if a.Articles[0].ID == 99 {
		t.Fatalf(`p[0].ID = %d, want = 99`, a.Articles[0].ID)
	}
}

func TestSearchReturnsNoArticleWhenNoMatchi(t *testing.T) {
	c := tests.Context()
	a, err := Search(c, Query{Keywords: "crazy world"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: ""crazy world"}) = %v, want nil`, err.Error())
	}

	if a.Total > 0 {
		t.Fatalf(`p.Total = %d, want == 0`, a.Total)
	}

	if len(a.Articles) > 0 {
		t.Fatalf(`len(p.Articles) = %d, want == 0`, len(a.Articles))
	}
}

func TestSearchReturnsArticlesWhenDescriptionIsFound(t *testing.T) {
	c := tests.Context()
	a, err := Search(c, Query{Keywords: "antiseptique"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: "antiseptique"}) = %v, want nil`, err.Error())
	}

	if a.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, a.Total)
	}

	if len(a.Articles) == 0 {
		t.Fatalf(`len(p.Articles) = %d, want > 0`, len(a.Articles))
	}

	if a.Articles[0].ID == 0 {
		t.Fatalf(`p[0].ID = %d, want > 0`, a.Articles[0].ID)
	}
}

func TestSearchReturnsArticleWhenSlugIsFound(t *testing.T) {
	c := tests.Context()
	slug := stringutil.Slugify(db.Escape(article.Title))
	a, err := Search(c, Query{Slug: slug}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Slug: slug}) = %v, want nil`, err.Error())
	}

	if a.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, a.Total)
	}

	if len(a.Articles) == 0 {
		t.Fatalf(`len(p.Articles) = %d, want > 0`, len(a.Articles))
	}

	if a.Articles[0].ID == 0 {
		t.Fatalf(`p[0].ID = %d, want > 0`, a.Articles[0].ID)
	}
}

func TestSearchReturnsNoArticleWhenCriteriaDoNotMatch(t *testing.T) {
	c := tests.Context()
	a, err := Search(c, Query{Keywords: "crazy"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: "crazy"}) = %v, want nil`, err.Error())
	}

	if a.Total != 0 {
		t.Fatalf(`p.Total = %d, want 0`, a.Total)
	}

	if len(a.Articles) != 0 {
		t.Fatalf(`len(p.Articles) = %d, want 0`, len(a.Articles))
	}
}

func TestDeleteReturnNilSuccess(t *testing.T) {
	c := tests.Context()

	image := path.Join(conf.ImgProxy.Path, "blog", "98")
	os.Create(image)

	if err := Delete(c, 98); err != nil {
		t.Fatalf(`a.Delete(c, 98) = %v, want nil`, err.Error())
	}
}

func TestFindReturnsErrorWhenIdsMissing(t *testing.T) {
	c := tests.Context()
	if _, err := Find(c, 0); err == nil || err.Error() != "input:id" {
		t.Fatalf(`Find(c,"") = %v, want "input:id"`, err.Error())
	}
}

func TestFindReturnsArticleWhenSuccess(t *testing.T) {
	c := tests.Context()
	p, err := Find(c, 99)
	if err != nil {
		t.Fatalf(`Find(c, "99") = %v, want nil`, err.Error())
	}

	if p.Title == "" {
		t.Fatalf(`p.Title = %v, want string`, p.Title)
	}

	if p.Description == "" {
		t.Fatalf(`p.Description  = %v, want string`, p.Description)
	}

	if p.Status != "online" {
		t.Fatalf(`p.Status = %v, want string`, p.Status)
	}

	if p.Image == "" {
		t.Fatalf(`p.Image = %v, want string`, p.Image)
	}
}

func TestDeletableReturnsFalseWhenTypeIsCms(t *testing.T) {
	c := tests.Context()

	if deletable, err := Deletable(c, 97); deletable || err != nil {
		t.Fatalf(`Deletable(c, 97) = %v, %v, want false', nil`, deletable, err)
	}
}

func TestDeletableReturnsFalseWhenTypeIsBlog(t *testing.T) {
	c := tests.Context()

	if deletable, err := Deletable(c, 99); !deletable || err != nil {
		t.Fatalf(`Deletable(c, c) = %v, %v, want true', nil`, !deletable, err)
	}
}
