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
	Status:      "online",
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

// func TestSaveReturnsErrorWhenImageIsEmpty(t *testing.T) {
// 	c := tests.Context()

// 	a := article
// 	a.Image = ""

// 	if err := a.Save(c); err == nil || err.Error() != "input_image_required" {
// 		t.Fatalf(`a.Save(c) = %v, want "input_image_required"`, err.Error())
// 	}
// }

func TestSaveReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	os.Create(path.Join(conf.WorkingSpace, "web", "tmp", "hello"))

	if err := article.Save(c); err != nil {
		t.Fatalf(`a.Save(c) = %v, want nil`, err.Error())
	}
}

// func TestListReturnsArticlesWhenSuccess(t *testing.T) {
// 	c := tests.Context()
// 	page := 0

// 	articles, err := List(c, page)
// 	if err != nil {
// 		t.Fatalf(`List(c, page) = %v, want nil`, err.Error())
// 	}

// 	a := articles[len(articles)-1]

// 	if a.ID == 0 {
// 		t.Fatalf(`a.ID = %v, want positive`, a.ID)
// 	}

// 	if a.Title == "" {
// 		t.Fatalf(`a.Title = %v, want not empty`, a.Title)
// 	}

// 	if a.Slug == "" {
// 		t.Fatalf(`a.Slug = %v, want not empty`, a.Slug)
// 	}

// 	if a.Description == "" {
// 		t.Fatalf(`a.Description = %v, want not empty`, a.Description)
// 	}

//		if a.Image == "" {
//			t.Fatalf(`a.Image = %v, want not empty`, a.Image)
//		}
//	}

func TestSearchReturnsArticlesWhenTitleIsFound(t *testing.T) {
	c := tests.Context()
	a, err := Search(c, Query{Keywords: "Manger de l'ail"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: "Manger de l'ail"}) = %v, want nil`, err.Error())
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

func TestSearchReturnsNoArticleWhenCriteriaDoNotMatch(t *testing.T) {
	c := tests.Context()
	a, err := Search(c, Query{Keywords: "hello"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: "hello"}) = %v, want nil`, err.Error())
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

	image := path.Join(conf.ImgProxy.Path, "blog", "3")
	os.Create(image)

	if err := Delete(c, 3); err != nil {
		t.Fatalf(`a.Delete(c, 3) = %v, want nil`, err.Error())
	}
}

func TestFindReturnsErrorWhenIdsMissing(t *testing.T) {
	c := tests.Context()
	if _, err := Find(c, 0); err == nil || err.Error() != "input_id_required" {
		t.Fatalf(`Find(c,"") = %v, want "input_id_required"`, err.Error())
	}
}

func TestFindReturnsArticleWhenSuccess(t *testing.T) {
	c := tests.Context()
	p, err := Find(c, 1)
	if err != nil {
		t.Fatalf(`Find(c, "PDT1") = %v, want nil`, err.Error())
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
