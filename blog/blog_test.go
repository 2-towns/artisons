package blog

import (
	"artisons/db"
	"artisons/tests"
	"errors"
	"fmt"
	"path"
	"reflect"
	"runtime"
	"testing"
)

var article Article = Article{
	Title:       "Mangez de l'ail !",
	Description: "C'est un antiseptique.",
	Slug:        db.Escape("Mangez de l'ail !"),
	Status:      "online",
}

var cur string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	cur = path.Dir(filename) + "/"
}

func TestValidate(t *testing.T) {
	ctx := tests.Context()

	var tests = []struct {
		name  string
		field string
		value string
		err   error
	}{
		{"title", "Title", "", errors.New("input:title")},
		{"description", "Description", "", errors.New("input:description")},
		{"slug", "Slug", "", errors.New("input:slug")},
		{"status", "Status", "", errors.New("input:status")},
		{"success", "", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := article

			if tt.field != "" {
				reflect.ValueOf(&a).Elem().FieldByName(tt.field).SetString(tt.value)
			}

			if err := a.Validate(ctx); fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want %s`, err, tt.err)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	ctx := tests.Context()

	tests.Del(ctx, "blog")
	tests.ImportData(ctx, cur+"testdata/article.redis")

	var tests = []struct {
		name  string
		query Query
		count int
	}{
		{"by title", Query{Keywords: "ail"}, 1},
		{"by title with space", Query{Keywords: "hello ail"}, 1},
		{"by description", Query{Keywords: "antiseptique"}, 1},
		{"by slug", Query{Slug: "mangez-de-lail"}, 1},
		{"by keywords", Query{Slug: "crazy"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := Search(ctx, tt.query, 0, 10)
			if err != nil {
				t.Fatalf(`err = %v, want nil`, err.Error())
			}

			if a.Total != tt.count {
				t.Fatalf(`p.Total = %d, want %d`, a.Total, tt.count)
			}

			if len(a.Articles) != tt.count {
				t.Fatalf(`len(p.Articles) = %d, want > %d`, len(a.Articles), tt.count)
			}
		})
	}

}

func TestDelete(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/article.redis")

	if err := Delete(ctx, 1); err != nil {
		t.Fatalf(`err = %v, want nil`, err.Error())
	}
}

func TestFind(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/article.redis")

	t.Run("id=0", func(t *testing.T) {
		if _, err := Find(ctx, 0); err == nil || err.Error() != "input:id" {
			t.Fatalf(`err = %v, want input:id`, err.Error())
		}
	})

	t.Run("id=1", func(t *testing.T) {
		c := tests.Context()
		p, err := Find(c, 1)
		if err != nil {
			t.Fatalf(`err = %v, want nil`, err.Error())
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
	})
}

func TestDeletable(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/article.redis")

	t.Run("type=cms", func(t *testing.T) {
		if deletable, err := Deletable(ctx, 2); deletable || err != nil {
			t.Fatalf(`deletable, err = %v, %v, want false, nil`, deletable, err)
		}
	})

	t.Run("type=blog", func(t *testing.T) {
		if deletable, err := Deletable(ctx, 1); !deletable || err != nil {
			t.Fatalf(`deletable, err = %v, %v, want true, nil`, !deletable, err)
		}
	})
}

func TestSave(t *testing.T) {
	ctx := tests.Context()

	t.Run("id=''", func(t *testing.T) {
		a := article
		if id, err := a.Save(ctx); id == "" || err != nil {
			t.Fatalf(`id, err = %s, %v, want not empty, nil`, id, err)
		}
	})

	t.Run("id=1", func(t *testing.T) {
		a := article
		a.ID = 1

		if id, err := a.Save(ctx); id != "1" || err != nil {
			t.Fatalf(`id, err = %s, %v, want 1, nil`, id, err)
		}
	})
}
