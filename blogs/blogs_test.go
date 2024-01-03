package blogs

import (
	"gifthub/conf"
	"gifthub/tests"
	"os"
	"path"
	"testing"
)

var article Article = Article{
	Title:       "La palestine doit être sauvée",
	Description: "On attend une réponse des pays musulmans.",
	Image:       path.Join(conf.WorkingSpace, "web", "tmp", "hello"),
}

func TestSaveReturnsErrorWhenTitleIsEmpty(t *testing.T) {
	c := tests.Context()

	a := article
	a.Title = ""

	if err := a.Save(c); err == nil || err.Error() != "input_title_required" {
		t.Fatalf(`a.Save(c) = %v, want "input_title_required"`, err.Error())
	}
}

func TestSaveReturnsErrorWhenDescriptionIsEmpty(t *testing.T) {
	c := tests.Context()

	a := article
	a.Description = ""

	if err := a.Save(c); err == nil || err.Error() != "input_description_required" {
		t.Fatalf(`a.Save(c) = %v, want "input_description_required"`, err.Error())
	}
}

func TestSaveReturnsErrorWhenImageIsEmpty(t *testing.T) {
	c := tests.Context()

	a := article
	a.Image = ""

	if err := a.Save(c); err == nil || err.Error() != "input_image_required" {
		t.Fatalf(`a.Save(c) = %v, want "input_image_required"`, err.Error())
	}
}

func TestSaveReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	os.Create(path.Join(conf.WorkingSpace, "web", "tmp", "hello"))

	if err := article.Save(c); err != nil {
		t.Fatalf(`a.Save(c) = %v, want nil`, err.Error())
	}
}

func TestListReturnsArticlesWhenSuccess(t *testing.T) {
	c := tests.Context()
	page := 0

	articles, err := List(c, page)
	if err != nil {
		t.Fatalf(`List(c, page) = %v, want nil`, err.Error())
	}

	a := articles[len(articles)-1]

	if a.ID == 0 {
		t.Fatalf(`a.ID = %v, want positive`, a.ID)
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

func TestDeleteReturnNilSuccess(t *testing.T) {
	c := tests.Context()

	image := path.Join(conf.ImgProxy.Path, "blog", "3")
	os.Create(image)

	if err := Delete(c, 3); err != nil {
		t.Fatalf(`a.Delete(c, 3) = %v, want nil`, err.Error())
	}
}
