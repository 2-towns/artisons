package blogs

import (
	"fmt"
	"gifthub/conf"
	"gifthub/tests"
	"os"
	"testing"
)

var article Article = Article{
	Title:       "La palestine doit être sauvée",
	Description: "On attend une réponse des pays musulmans.",
	Image:       "/tmp/hello",
}

func TestSaveWithoutTitle(t *testing.T) {
	c := tests.Context()

	a := article
	a.Title = ""

	if err := a.Save(c); err == nil || err.Error() != "article_title_required" {
		t.Fatalf(`a.Save(c) = %v, want "article_title_required"`, err.Error())
	}
}

func TestSaveWithoutDescription(t *testing.T) {
	c := tests.Context()

	a := article
	a.Description = ""

	if err := a.Save(c); err == nil || err.Error() != "article_description_required" {
		t.Fatalf(`a.Save(c) = %v, want "article_description_required"`, err.Error())
	}
}

func TestSaveWithoutImage(t *testing.T) {
	c := tests.Context()

	a := article
	a.Image = ""

	if err := a.Save(c); err == nil || err.Error() != "article_image_required" {
		t.Fatalf(`a.Save(c) = %v, want "article_image_required"`, err.Error())
	}
}

func TestSave(t *testing.T) {
	c := tests.Context()

	os.Create("/tmp/hello")

	if err := article.Save(c); err != nil {
		t.Fatalf(`a.Save(c) = %v, want nil`, err.Error())
	}
}

func TestList(t *testing.T) {
	c := tests.Context()
	page := 0

	articles, err := List(c, page)
	if err != nil {
		t.Fatalf(`List(c, page) = %v, want nil`, err.Error())
	}

	if len(articles) == 0 {
		t.Fatalf(`len(articles) = %v, want 1`, len(articles))
	}

	a := articles[0]

	if a.ID != 1 {
		t.Fatalf(`a.ID = %v, want 1`, a.ID)
	}

	if a.Title == "" {
		t.Fatalf(`a.Title = %v, want not empty`, a.Title)
	}

	if a.Slug == "" {
		t.Fatalf(`a.Slug = %v, want not empty`, a.Slug)
	}

	if a.Description == "" {
		t.Fatalf(`a.Description = %v, want not empty`, a.Description)
	}

	if a.Image == "" {
		t.Fatalf(`a.Image = %v, want not empty`, a.Image)
	}
}

func TestDelete(t *testing.T) {
	c := tests.Context()

	image := fmt.Sprintf("%s/articles/%d", conf.ImgProxyPath, 3)
	os.Create(image)

	if err := Delete(c, 3); err != nil {
		t.Fatalf(`a.Delete(c, 3) = %v, want nil`, err.Error())
	}
}
