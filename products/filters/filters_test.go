package filters

import (
	"gifthub/db"
	"gifthub/tests"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

var filter Filter = Filter{
	Key:       "color",
	Type:      "color",
	Label:     "color",
	Score:     1,
	Values:    []string{"yellow", "blue"},
	Active:    true,
	UpdatedAt: time.Now(),
}

func init() {
	ctx := tests.Context()

	db.Redis.HSet(ctx, "filter:colors",
		"key", filter.Key,
		"label", filter.Label,
		"values", strings.Join(filter.Values, ";"),
		"updated_at", filter.UpdatedAt.Unix(),
	)

	db.Redis.HSet(ctx, "filter:sizes",
		"key", "sizes",
		"label", "Sizes",
		"values", "S;M;L",
		"updated_at", time.Now().Unix(),
	)

	db.Redis.ZAdd(ctx, "filters", redis.Z{
		Score:  float64(filter.UpdatedAt.Unix()),
		Member: filter.Key,
	}, redis.Z{
		Score:  float64(filter.UpdatedAt.Unix()),
		Member: "sizes",
	})

	db.Redis.ZAdd(ctx, "fitlter:active", redis.Z{
		Score:  float64(filter.Score),
		Member: filter.Key,
	})
}

func TestValidateReturnsErrorWhenTheKeyIsEmpty(t *testing.T) {
	c := tests.Context()

	f := filter
	f.Key = ""

	if err := f.Validate(c); err == nil || err.Error() != "input:key" {
		t.Fatalf(`f.Validate(c) = %v, want 'input:key'`, err.Error())
	}
}

func TestValidateReturnsErrorWhenTypeIsEmpty(t *testing.T) {
	c := tests.Context()

	f := filter
	f.Type = ""

	if err := f.Validate(c); err == nil || err.Error() != "input:type" {
		t.Fatalf(`ta.Validate(c) = %v, want 'input:type'`, err.Error())
	}
}

func TestValidateReturnsErrorWhenTypeIsInvalid(t *testing.T) {
	c := tests.Context()

	f := filter
	f.Type = "hello"

	if err := f.Validate(c); err == nil || err.Error() != "input:type" {
		t.Fatalf(`ta.Validate(c) = %v, want 'input:type'`, err.Error())
	}
}

func TestValidateReturnsErrorWhenLabelIsEmpty(t *testing.T) {
	c := tests.Context()

	f := filter
	f.Label = ""

	if err := f.Validate(c); err == nil || err.Error() != "input:label" {
		t.Fatalf(`ta.Validate(c) = %v, want 'input:label'`, err.Error())
	}
}

func TestFindReturnsFilterWhenTheKeyExists(t *testing.T) {
	c := tests.Context()

	filter, err := Find(c, "colors")

	if err != nil {
		t.Fatalf(`Find(c, "colors") = %v, want nil`, err.Error())
	}

	if filter.Key == "" {
		t.Fatalf(`tag.Key = %s, want not empty`, filter.Key)
	}

	if filter.Label == "" {
		t.Fatalf(`tag.Label = %s, want not empty`, filter.Label)
	}
}

func TestFindReturnsEmptyFilterWhenTheKeyDoesNotExist(t *testing.T) {
	c := tests.Context()

	if _, err := Find(c, "hello"); err == nil || err.Error() != "oops the data is not found" {
		t.Fatalf(`Find(c, "hello") = %v, want nil`, err.Error())
	}
}

func TestSaveReturnsNilWhenEmptyWhenSuccess(t *testing.T) {
	c := tests.Context()

	if _, err := filter.Save(c); err != nil {
		t.Fatalf(`filter.Save(c) = %v, want nil`, err)
	}
}

func TestListReturnsFilters(t *testing.T) {
	c := tests.Context()

	r, err := List(c, 0, 10)
	if err != nil {
		t.Fatalf(`List(c) = %v, want nil`, err)
	}

	if r.Total == 0 {
		t.Fatalf(`r.Total = %d, want > 0`, r.Total)
	}

	if len(r.Filters) == 0 {
		t.Fatalf(`len(r.Filters) = %d, want > 0`, len(r.Filters))
	}

	filter := r.Filters[0]

	if filter.Key == "" {
		t.Fatalf(`filter.Key = %s, want not empty`, filter.Key)
	}
}

func TestDeleteReturnsErrorWhenKeyIsEmpty(t *testing.T) {
	c := tests.Context()

	if err := Delete(c, ""); err == nil || err.Error() != "input:key" {
		t.Fatalf(`Delete(c) = %v, want 'input:key'`, err)
	}
}

func TestDeleteReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := Delete(c, "sizes"); err != nil {
		t.Fatalf(`Delete(c) = %v, want nil`, err)
	}
}

func TestExistsReturnsTrueWhenFilterDoesNotExist(t *testing.T) {
	c := tests.Context()

	if exists, err := Exists(c, "colors"); exists == false || err != nil {
		t.Fatalf(`Exists(c, "colors") = %v, %v, want true', nil`, exists, err)
	}
}

func TestExistsReturnsFalseWhenFilterDoesNotExist(t *testing.T) {
	c := tests.Context()

	if exists, err := Exists(c, "hello"); exists == true || err != nil {
		t.Fatalf(`Exists(c, "hello") = %v, %v, want false', nil`, exists, err)
	}
}
